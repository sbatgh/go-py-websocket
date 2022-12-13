import sys
import os
from pathlib import Path

import argparse
from time import time, sleep
import json

import websockets
import pandas as pd

import asyncio


# get real time messages from websockets
async def get_realtime_messages(url_):
    async with websockets.connect(url_) as websocket:
        print('Websocket is connected')
        while True:
            if not websocket.open:
                print("Websocket is closed")
                sleep(10)
                websocket = await websockets.connect(url_)
                try:
                    # print('in try')
                    message = await websocket.recv()
                    msg = json.loads(message)
                    print(msg)
                except Exception as e:
                    print(e)
                    print("Websocket is closed")
                    sleep(10)
                    
            try:
                print('second try')
                message = await websocket.recv()
                msg = json.loads(message)
                print(msg)
            except Exception as e:
                print(e)
                print("Websocket is closed")
                sleep(10)

if __name__ == "__main__":

    url = "ws://localhost:8080/ws"
    loop = asyncio.get_event_loop()
    several_futures = asyncio.gather(get_realtime_messages(url))
    try:
        loop.run_until_complete(several_futures)
    except KeyboardInterrupt:
        print("KeyboardInterrupt")
        loop.close()
    finally:
        loop.close()


