import aiohttp
import asyncio
import json, time


async def is_prime(number):
    for i in range(2, int(number/2)):
        if number%i == 0:
            return False
    return True

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
        story_detail = json.loads(html)
        prime_result = await is_prime(story_detail['time'])
        print (f"Story has title {story_detail['title']} and ID is Prime - {prime_result}")

        

async def main():
    new_stories = await get_new_stories()
    worker_count = 25
    limit = len(new_stories)
    for i in range(limit):
        tasks = []
        for j in range(i, i+worker_count):
            if j == limit:
                break
            tasks.append(asyncio.create_task(get_story_detail(new_stories[j])))
        await asyncio.gather(*tasks)
        i = i+worker_count    

if __name__ == '__main__':
    start_time = time.time()
    print(start_time)
    asyncio.run(main())
    elapsed = time.time() - start_time
    print(f"Total seconds to finish - {elapsed}")