import asyncio
import logging
import os
from asyncio import Queue
from dataclasses import dataclass, asdict, field
from typing import Optional, List, Tuple

from aiohttp import InvalidURL
from yarl import URL

import aiohttp


FORMAT = '%(asctime)-15s %(name)-12s %(levelname)-8s %(message)s'
logging.basicConfig(level=logging.INFO, format=FORMAT)
logger = logging.getLogger(__name__)


@dataclass
class License:
    id: Optional[int] = None
    digAllowed: int = 0
    digUsed: int = 0


@dataclass
class Area:
    posX: int
    posY: int
    sizeX: int = 1
    sizeY: int = 1


@dataclass
class Dig:
    licenseID: int
    posX: int
    posY: int
    depth: int = 1


@dataclass(frozen=True, order=True)
class Explore:
    priority: int
    area: Area = field(compare=False)
    amount: int = field(compare=False)


@dataclass
class Wallet:
    balance: int = 0
    wallet: List[int] = field(default_factory=list)


@dataclass(order=True)
class Treasure:
    priority: int
    treasures: List[str] = field(default_factory=list, compare=False)


class Client:

    def __init__(self, base_url: URL) -> None:
        self.base_url = base_url
        self._client = aiohttp.ClientSession()
        license = License()
        self.license = license

    async def close(self) -> None:
        return await self._client.close()

    async def post_license(self, coin: List[int]) -> Optional[License]:
        logger.debug('get license')
        url = self.base_url / 'licenses'
        try:
            async with self._client.post(url, json=coin) as resp:
                if resp.status == 200:
                    _json = await resp.json()
                    return License(**_json)
        except InvalidURL:
            logger.error('invalid url %s', url)
            return

    async def get_license(self) -> Optional[List[License]]:
        logger.debug('get list license')
        url = self.base_url / 'licenses'
        try:
            async with self._client.get(url) as resp:
                if resp.status == 200:
                    _json = await resp.json()
                    license_list = [License(**item) for item in _json]
                    return license_list
        except InvalidURL:
            logger.error('invalid url %s', url)
            return

    async def post_dig(self, dig: Dig) -> Optional[Treasure]:
        logger.debug('dig')
        url = self.base_url / 'dig'
        try:
            async with self._client.post(url, json=asdict(dig)) as resp:
                if resp.status == 200:
                    treasures = await resp.json()
                    return Treasure(priority=0, treasures=treasures)
                elif resp.status == 403:
                    self.license.id = None
        except InvalidURL:
            logger.error('invalid url %s', url)
            return

    async def post_cash(self, treasure: str):
        logger.debug('cash')
        url = self.base_url / 'cash'
        try:
            async with self._client.post(url, json=treasure) as resp:
                if resp.status == 200:
                    return await resp.json()
        except InvalidURL:
            logger.error('invalid url %s', url)
            return

    async def post_explore(self, area: Area) -> Optional[Explore]:
        logger.debug('explore')
        url = self.base_url / 'explore'
        try:
            async with self._client.post(url, json=asdict(area)) as resp:
                if resp.status != 200:
                    return
                _json = await resp.json()
                area = Area(**_json['area'])
                return Explore(priority=0, area=area, amount=_json.get('amount', 0))
        except InvalidURL:
            logger.error('invalid url: %s', url)

    async def get_balance(self) -> Optional[Wallet]:
        logger.debug('get balance')
        url = self.base_url / 'balance'
        try:
            async with self._client.get(url) as resp:
                if resp.status != 200:
                    return
                _json = await resp.json()
                wallet = Wallet(**_json)
            return wallet
        except InvalidURL:
            logger.error('invalid url %s', url)
            return

    async def update_license(self) -> None:
        logger.debug('update license')
        coins = []
        license = await self.post_license(coins)
        if license is not None:
            self.license = license


async def game(client: Client):
    try:
        for x in range(3500):
            for y in range(3500):
                area = Area(x, y)
                result = await client.post_explore(area)
                if result is None or not result.amount:
                    continue
                depth = 1
                left = result.amount
                while depth <= 10 and left > 0:
                    while client.license.id is None or client.license.digUsed >= client.license.digAllowed:
                        await client.update_license()

                    dig = Dig(licenseID=client.license.id, posX=x, posY=y, depth=depth)
                    treasures = await client.post_dig(dig)
                    client.license.digUsed += 1
                    depth += 1
                    if treasures is not None:
                        for treasure in treasures.treasures:
                            res = await client.post_cash(treasure)
                            if res:
                                left -= 1
    except Exception as e:
        logger.error('error: %s', e)
    finally:
        await client.close()


async def main():
    logger.debug('start')
    address = os.getenv('ADDRESS')
    base_url = URL.build(scheme='http', host=address, port=8000)
    logger.debug('base url: %s', base_url)
    await game(client=Client(base_url))


if __name__ == '__main__':
    asyncio.run(main())
