# Usage: python get_rss_url.py
# Exported links from OneTab can be used as is.
# Save it as url.txt and place on the same directory as get_rss_url.py

# Required packages: requests, BeautifulSoup4

import requests
from bs4 import BeautifulSoup

f = open('url.txt', 'r')
data = f.read()
f.close()

urls = data.split('\n')

for raw_url in urls:
    if not raw_url.startswith("https://"):
        continue
    url = raw_url.split(" | ", 1)[0]

    source = requests.get(url).text
    soup = BeautifulSoup(source, features="html.parser")

    rss_url = soup.find("link", attrs={'title': 'RSS'})
    channel_name = soup.find("meta", property="og:title", content=True)
    if rss_url is None or channel_name is None:
        continue

    print("# " + channel_name["content"])
    print(rss_url["href"])
