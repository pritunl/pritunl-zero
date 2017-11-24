#!/usr/bin/env python
import os
import json
import urllib2
import subprocess
import urlparse
import sys

SSH_DIR = '~/.ssh'
CONF_PATH = SSH_DIR + '/pritunl-zero.json'
VERSION = '1.0.731.28'

USAGE = """\
Usage: pritunl-ssh [command]

Commands:
  help      Show help
  version   Print the version and exit
  config    Reconfigure options"""

zero_server = None
pub_key_path = None
ssh_dir_path = os.path.expanduser(SSH_DIR)
conf_path = os.path.expanduser(CONF_PATH)
changed = False

if '--help' in sys.argv[1:] or 'help' in sys.argv[1:]:
    print USAGE
    exit()

if '--version' in sys.argv[1:] or 'version' in sys.argv[1:]:
    print 'pritunl-ssh v' + VERSION
    exit()

if '--config' not in sys.argv[1:] and \
        'config' not in sys.argv[1:] and \
    os.path.isfile(conf_path):
    with open(conf_path, 'r') as conf_file:
        conf_data = conf_file.read()
        try:
            conf_data = json.loads(conf_data)
            zero_server = conf_data.get('server')
            pub_key_path = conf_data.get('public_key_path')
        except:
            print 'WARNING: Failed to parse config file'

if not zero_server:
    server = raw_input('Enter Pritunl Zero user hostname: ')
    server_url = urlparse.urlparse(server)
    zero_server = 'https://%s' % (server_url.netloc or server_url.path)
    changed = True

print 'SERVER: ' + zero_server

if not pub_key_path or not os.path.exists(os.path.expanduser(pub_key_path)):
    if not os.path.exists(ssh_dir_path):
        print 'ERROR: No SSH keys found, run "ssh-keygen" to create a key'
        exit()

    ssh_names = []

    print 'Select SSH key:'

    for filename in os.listdir(ssh_dir_path):
        if '.pub' not in filename:
            continue

        ssh_names.append(filename)
        print '[%d] %s' % (len(ssh_names), filename)

    key_input = raw_input('Enter key number or full path to key: ')

    try:
        index = int(key_input)
        pub_key_path = os.path.join(SSH_DIR, ssh_names[index - 1])
    except ValueError, IndexError:
        pass

    if not pub_key_path:
        if key_input in ssh_names:
            pub_key_path = os.path.join(SSH_DIR, key_input)
        else:
            pub_key_path = key_input

    pub_key_path = os.path.normpath(pub_key_path)
    changed = True

pub_key_path_full = os.path.expanduser(pub_key_path)
cert_path = pub_key_path.rsplit('.pub', 1)[0] + '-cert.pub'
cert_path_full = os.path.expanduser(cert_path)
if not os.path.exists(pub_key_path_full):
    print 'ERROR: Selected SSH key does not exist'
    exit()

if not pub_key_path_full.endswith('.pub'):
    print 'ERROR: SSH key path must end with .pub'
    exit()

print 'SSH_KEY: ' + pub_key_path

with open(conf_path, 'w') as conf_file:
    conf_file.write(json.dumps({
        'server': zero_server,
        'public_key_path': pub_key_path,
    }))

with open(pub_key_path_full, 'r') as pub_key_file:
    pub_key_data = pub_key_file.read().strip()

req = urllib2.Request(
    zero_server + '/ssh/challenge',
    data=json.dumps({
        'public_key': pub_key_data,
    }),
)
req.add_header('Content-Type', 'application/json')
req.get_method = lambda: 'POST'
try:
    resp = urllib2.urlopen(req)
    resp_data = resp.read()
    status_code = resp.getcode()
except urllib2.HTTPError as exception:
    status_code = exception.code
    resp_data = ''

if status_code != 200:
    print 'ERROR: SSH challenge request failed with status %d' % status_code
    if resp_data:
        print resp_data
    exit()

token = json.loads(resp_data)['token']

token_url = zero_server + '/ssh?ssh-token=' + token

print 'OPEN: ' + token_url

subprocess.Popen(['open', token_url])

for i in xrange(3):
    req = urllib2.Request(
        zero_server + '/ssh/challenge',
        data=json.dumps({
            'public_key': pub_key_data,
            'token': token,
        }),
    )
    req.add_header('Content-Type', 'application/json')
    req.get_method = lambda: 'PUT'

    try:
        resp = urllib2.urlopen(req)
        status_code = resp.getcode()
        resp_data = resp.read()
    except urllib2.HTTPError as exception:
        status_code = exception.code
        resp_data = ''

    if status_code == 205:
        continue
    break

if status_code == 205:
    print 'ERROR: SSH verification request timed out'
    exit()

if status_code == 401:
    print 'ERROR: SSH verification request was denied'
    exit()

if status_code != 200:
    print 'ERROR: SSH verification failed with status %d' % status_code
    if resp_data:
        print resp_data
    exit()

certificates = json.loads(resp_data)['certificates']

with open(cert_path_full, 'w') as cert_file:
    cert_file.write('\n'.join(certificates) + '\n')

print 'CERTIFICATE: ' + cert_path
print 'Successfully validated SSH key'
