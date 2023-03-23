
import random
import string
import subprocess
from time import process_time
import time

urls = []
TOTAL_REQUESTS = 2500

# NOTE:  COPIED STRUCTURE FROM A1 code

def GenRandomURLPair():
    longResource = "http://"+''.join(random.choice(string.ascii_uppercase + string.digits) for _ in range(10))+'.com'
    shortResource = ''.join(random.choice(string.ascii_uppercase + string.digits) for _ in range(6))
    return shortResource, longResource

def put():
    shortResource, longResource = GenRandomURLPair()
    request="http://localhost:4000/?short={}&long={}".format(shortResource, longResource)
    if subprocess.call(["curl", "-4", "-X", "PUT", request], stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL) == 0:
        urls.append((shortResource, longResource))
        return True
    return False

def get():
    shortResource, _ = random.choice(urls)
    request="http://localhost:4000/{}".format(shortResource)
    if subprocess.call(["curl", "-4", request], stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL) == 0:
        return True
    return False

if __name__ == "__main__":
    successfulGets = 0
    successfulPosts = 0
    t1 = time.time()
    for _ in range(TOTAL_REQUESTS):
        successfulPosts += int(put())
    post_time = time.time() - t1

    t1 = time.time()
    for _ in range(TOTAL_REQUESTS):
        successfulGets += int(get())
    get_time = time.time() - t1
    print("{} Requests -> PUT: {}, GET: {}".format(TOTAL_REQUESTS, post_time, get_time))
    