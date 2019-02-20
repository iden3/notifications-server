#!/usr/bin/env python3
"""notifications-server endpoints test
"""

import requests
import provoj

URL = "http://127.0.0.1:10000/api/unstable"

t = provoj.NewTest("notifications-server")

r = requests.get(URL + "/")
t.rStatus("get info", r)

idaddr0 = "0x47a2b2353f1a55e4c975b742a7323c027160b4e3"
idaddr1 = "0xd9d6800a1b20ceebef5420f878bbd915f8b4ed85"

r = requests.post(URL + "/notifications/" + idaddr0, json={"data": "notif00"})
t.rStatus("post notification to " + idaddr0, r)

for i in range(10):
    notificationData = "notif" + str(i)
    r = requests.post(URL + "/notifications/" + idaddr0, json={"data": notificationData})
    r = requests.post(URL + "/notifications/" + idaddr1, json={"data": notificationData})

r = requests.post(URL + "/notifications", json={"idAddr": idaddr0})
t.rStatus("get notifications for " + idaddr0, r)
jsonR = r.json()

r = requests.delete(URL + "/notifications", json={"idAddr": idaddr0})
t.rStatus("delete notifications for " + idaddr0, r)
jsonR = r.json()

r = requests.post(URL + "/notifications", json={"idAddr": idaddr0})
t.rStatus("get notifications for " + idaddr0, r)
jsonR = r.json()
t.equal("expect no notifications for idaddr0", jsonR["notifications"], None)

r = requests.post(URL + "/notifications", json={"idAddr": idaddr1})
t.rStatus("get notifications for " + idaddr1, r)
jsonR = r.json()
t.notequal("expect notifications for idaddr1", len(jsonR["notifications"]), 0)

t.printScores()
