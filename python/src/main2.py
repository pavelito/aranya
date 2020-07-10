import aiohttp
import asyncio
import json, time

async def is_prime(number):
    flag = True
    for i in range(2, int(number/2)):
        if number % i == 0:
            flag = False
            break
    return flag

async def countdown(number):
    counter = number
    while counter > 0:
        counter = counter - 1
    print (f"Countdown done for {number}")

async def fetch(session, url):
    async with session.get(url) as response:
        return await response.text()

async def get_new_stories():
    async with aiohttp.ClientSession() as session:
        html = await fetch(session, 'https://hacker-news.firebaseio.com/v0/topstories.json')
        return json.loads(html)

async def get_story_detail(story_id):
    async with aiohttp.ClientSession() as session:
        html = await fetch(session, f"https://hacker-news.firebaseio.com/v0/item/{story_id}.json?print=pretty")
        return json.loads(html)

async def main():
    new_stories = await get_new_stories()
    for story_id in new_stories:
        story_detail = await get_story_detail(story_id)
        print (f"Story has title {story_detail['title']}")
        await countdown(story_detail['time'])

    

if __name__ == '__main__':
    start_time = time.time()
    #print(start_time)
    loop = asyncio.get_event_loop()
    loop.run_until_complete(main())
    elapsed = time.time() - start_time
    print(f"Total seconds to finish - {elapsed}")