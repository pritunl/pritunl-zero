import requests
import time
import uuid
import hmac
import hashlib
import base64
import json

BASE_URL = 'https://'
API_TOKEN = ''
API_SECRET = ''

def auth_request(method, path, headers=None, data=None):
    auth_timestamp = str(int(time.time()))
    auth_nonce = uuid.uuid4().hex
    auth_string = '&'.join([API_TOKEN, auth_timestamp, auth_nonce,
        method.upper(), path])
    auth_signature = base64.b64encode(hmac.new(
        API_SECRET, auth_string, hashlib.sha512).digest())
    auth_headers = {
        'Pritunl-Zero-Token': API_TOKEN,
        'Pritunl-Zero-Timestamp': auth_timestamp,
        'Pritunl-Zero-Nonce': auth_nonce,
        'Pritunl-Zero-Signature': auth_signature,
    }
    if headers:
        auth_headers.update(headers)
    return getattr(requests, method.lower())(
        BASE_URL + path,
        headers=auth_headers,
        data=data,
    )

users = []

users.append({
    'type': 'local',
    'username': 'user1',
    'password': 'password1',
    'roles': ['role1', 'role2'],
})

for user in users:
    response = auth_request(
        'POST',
        '/user',
        headers={
            'Content-Type': 'application/json',
        },
        data=json.dumps(user),
    )

    assert(response.status_code == 200)
