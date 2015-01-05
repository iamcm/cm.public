from pprint import pprint
import requests
import json

data = {
    'Name': 'cmtest2',
    'IPs':['192.168.1.25']
}

strdata = json.dumps(data)

r = requests.post('http://localhost:8050/data', data={'data':strdata})

#pprint(r.json())
pprint(r.content)

