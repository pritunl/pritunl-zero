import optparse
import datetime
import re
import sys
import subprocess
import time
import math
import json
import requests
import os
import getpass
import base64
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.backends import default_backend
from cryptography.hazmat.primitives.ciphers import (
    Cipher, algorithms, modes
)
from cryptography.hazmat.primitives.kdf.pbkdf2 import PBKDF2HMAC

USAGE = """Usage: builder [command] [options]
Command Help: builder [command] --help

Commands:
  version               Print the version and exit
  set-version           Set current version
  build                 Build and release"""

CONSTANTS_PATH = 'constants/constants.go'
CHANGES_PATH = 'CHANGES'
STABLE_PACUR_PATH = '../pritunl-pacur'
TEST_PACUR_PATH = '../pritunl-pacur-test'
BUILD_KEYS_PATH = 'build_keys.json'
BUILD_TARGETS = ('pritunl-zero',)
REPO_NAME = 'pritunl-zero'

cur_date = datetime.datetime.utcnow()
pacur_path = None

def wget(url, cwd=None, output=None):
    if output:
        args = ['wget', '-O', output, url]
    else:
        args = ['wget', url]
    subprocess.check_call(args, cwd=cwd)

def post_git_asset(release_id, file_name, file_path):
    file_size = os.path.getsize(file_path)
    response = requests.post(
        'https://uploads.github.com/repos/%s/%s/releases/%s/assets' % (
            github_owner, REPO_NAME, release_id),
        verify=False,
        headers={
            'Authorization': 'token %s' % github_token,
            'Content-Type': 'application/octet-stream',
            'Content-Size': str(file_size),
        },
        params={
            'name': file_name,
        },
        data=open(file_path, 'rb').read(),
    )

    if response.status_code != 201:
        print('Failed to create asset on github')
        print(response.json())
        sys.exit(1)

def get_ver(version):
    day_num = (cur_date - datetime.datetime(2015, 11, 24)).days
    min_num = int(math.floor(((cur_date.hour * 60) + cur_date.minute) / 14.4))
    ver = re.findall(r'\d+', version)
    ver_str = '.'.join((ver[0], ver[1], str(day_num), str(min_num)))
    ver_str += ''.join(re.findall('[a-z]+', version))

    return ver_str

def get_int_ver(version):
    ver = re.findall(r'\d+', version)

    if 'snapshot' in version:
        pass
    elif 'alpha' in version:
        ver[-1] = str(int(ver[-1]) + 1000)
    elif 'beta' in version:
        ver[-1] = str(int(ver[-1]) + 2000)
    elif 'rc' in version:
        ver[-1] = str(int(ver[-1]) + 3000)
    else:
        ver[-1] = str(int(ver[-1]) + 4000)

    return int(''.join([x.zfill(4) for x in ver]))

def iter_packages():
    for target in BUILD_TARGETS:
        target_path = os.path.join(pacur_path, target)
        for name in os.listdir(target_path):
            if cur_version not in name:
                continue
            elif name.endswith(".pkg.tar.xz"):
                pass
            elif name.endswith(".rpm"):
                pass
            elif name.endswith(".deb"):
                pass
            else:
                continue

            path = os.path.join(target_path, name)

            yield name, path

# Parse args
if len(sys.argv) > 1:
    cmd = sys.argv[1]
else:
    cmd = 'version'

def aes_encrypt(passphrase, data):
    enc_salt = os.urandom(32)
    enc_iv = os.urandom(16)

    kdf = PBKDF2HMAC(
        algorithm=hashes.SHA1(),
        length=32,
        salt=enc_salt,
        iterations=1000,
        backend=default_backend(),
    )
    enc_key = kdf.derive(passphrase.encode())

    data += '\x00' * (16 - (len(data) % 16))

    cipher = Cipher(
        algorithms.AES(enc_key),
        modes.CBC(enc_iv),
        backend=default_backend()
    ).encryptor()
    enc_data = cipher.update(data.encode()) + cipher.finalize()

    return '\n'.join([
        base64.b64encode(enc_salt).decode('utf-8'),
        base64.b64encode(enc_iv).decode('utf-8'),
        base64.b64encode(enc_data).decode('utf-8'),
    ])

def aes_decrypt(passphrase, data):
    data = data.split('\n')
    if len(data) < 3:
        raise ValueError('Invalid encryption data')

    enc_salt = base64.b64decode(data[0])
    enc_iv = base64.b64decode(data[1])
    enc_data = base64.b64decode(data[2])

    kdf = PBKDF2HMAC(
        algorithm=hashes.SHA1(),
        length=32,
        salt=enc_salt,
        iterations=1000,
        backend=default_backend(),
    )
    enc_key = kdf.derive(passphrase.encode())

    cipher = Cipher(
        algorithms.AES(enc_key),
        modes.CBC(enc_iv),
        backend=default_backend()
    ).decryptor()
    data = cipher.update(enc_data) + cipher.finalize()

    return data.decode('utf-8').replace('\x00', '')

passphrase = getpass.getpass('Enter passphrase: ')

if cmd == 'encrypt':
    passphrase2 = getpass.getpass('Enter passphrase: ')

    if passphrase != passphrase2:
        print('ERROR: Passphrase mismatch')
        sys.exit(1)

    with open(BUILD_KEYS_PATH, 'r') as build_keys_file:
        data = build_keys_file.read().strip()

    enc_data = aes_encrypt(passphrase, data)

    with open(BUILD_KEYS_PATH, 'w') as build_keys_file:
        build_keys_file.write(enc_data)

    sys.exit(0)

if cmd == 'decrypt':
    with open(BUILD_KEYS_PATH, 'r') as build_keys_file:
        enc_data = build_keys_file.read().strip()

    data = aes_decrypt(passphrase, enc_data)

    with open(BUILD_KEYS_PATH, 'w') as build_keys_file:
        build_keys_file.write(data)

    sys.exit(0)

# Load build keys
with open(BUILD_KEYS_PATH, 'r') as build_keys_file:
    enc_data = build_keys_file.read()
    data = aes_decrypt(passphrase, enc_data)
    build_keys = json.loads(data.strip())
    github_owner = build_keys['github_owner']
    github_token = build_keys['github_token']
    gitlab_host = build_keys['gitlab_host']
    gitlab_token = build_keys['gitlab_token']
    mirror_url = build_keys['mirror_url']
    test_mirror_url = build_keys['test_mirror_url']


# Get package info
with open(CONSTANTS_PATH, 'r') as constants_file:
    cur_version = re.findall('= "(.*?)"', constants_file.read())[0]


parser = optparse.OptionParser(usage=USAGE)
(options, args) = parser.parse_args()

build_num = 0


# Run cmd
if cmd == 'version':
    print('%s v%s' % (REPO_NAME, cur_version))
    sys.exit(0)


if cmd == 'sync-releases':
    next_url = 'https://api.github.com/repos/%s/%s/releases' % (
        github_owner, REPO_NAME)

    while True:
        # Get github release
        response = requests.get(
            next_url,
            headers={
                'Authorization': 'token %s' % github_token,
                'Content-type': 'application/json',
            },
        )

        if response.status_code != 200:
            print('Failed to get repo releases on github')
            print(response.json())
            sys.exit(1)

        for release in response.json():
            print(release['tag_name'])

            # Create gitlab release
            resp = requests.post(
                'https://' + gitlab_host + '/api/v4/projects' + \
                '/%s%%2F%s/repository/tags/%s/release' % (
                    github_owner, REPO_NAME, release['tag_name']),
                headers={
                    'Private-Token': gitlab_token,
                    'Content-type': 'application/json',
                },
                data=json.dumps({
                    'tag_name': release['tag_name'],
                    'description': release['body'],
                }),
            )

            if resp.status_code not in (201, 409):
                print('Failed to create releases on gitlab')
                print(resp.json())
                sys.exit(1)

        if 'Link' not in response.headers or \
                'rel="next"' not in response.headers['Link']:
            break
        next_url = response.headers['Link'].split(';')[0][1:-1]


if cmd == 'set-version':
    new_version_orig = args[1]
    new_version = get_ver(new_version_orig)
    is_snapshot = 'snapshot' in new_version
    pacur_path = TEST_PACUR_PATH if is_snapshot else STABLE_PACUR_PATH


    # Update changes
    if not is_snapshot:
        with open(CHANGES_PATH, 'r') as changes_file:
            changes_data = changes_file.read()

        with open(CHANGES_PATH, 'w') as changes_file:
            ver_date_str = 'Version ' + new_version.replace(
                'v', '') + cur_date.strftime(' %Y-%m-%d')
            changes_file.write(changes_data.replace(
                '<%= version %>',
                '%s\n%s' % (ver_date_str, '-' * len(ver_date_str)),
            ))


    # Check for duplicate version
    response = requests.get(
        'https://api.github.com/repos/%s/%s/releases' % (
            github_owner, REPO_NAME),
        headers={
            'Authorization': 'token %s' % github_token,
            'Content-type': 'application/json',
        },
    )

    if response.status_code != 200:
        print('Failed to get repo releases on github')
        print(response.json())
        sys.exit(1)

    for release in response.json():
        if release['tag_name'] == new_version:
            print('Version already exists in github')
            sys.exit(1)


    # Generate changelog
    version = None
    release_body = ''
    if not is_snapshot:
        with open(CHANGES_PATH, 'r') as changelog_file:
            for line in changelog_file.readlines()[2:]:
                line = line.strip()

                if not line or line[0] == '-':
                    continue

                if line[:7] == 'Version':
                    if version:
                        break
                    version = line.split(' ')[1]
                elif version:
                    release_body += '* %s\n' % line

    if not is_snapshot and version != new_version:
        print('New version does not exist in changes')
        sys.exit(1)

    if is_snapshot:
        release_body = '* Snapshot release'
    elif not release_body:
        print('Failed to generate github release body')
        sys.exit(1)
    release_body = release_body.rstrip('\n')


    # Update constants
    with open(CONSTANTS_PATH, 'r') as constants_file:
        constants_data = constants_file.read()

    with open(CONSTANTS_PATH, 'w') as constants_file:
        constants_file.write(re.sub(
            '(= ".*?")',
            '= "%s"' % new_version,
            constants_data,
            count=1,
            ))

    # Git commit
    subprocess.check_call(['git', 'reset', 'HEAD', '.'])
    subprocess.check_call(['git', 'add', CHANGES_PATH])
    subprocess.check_call(['git', 'add', CONSTANTS_PATH])
    subprocess.check_call(['git', 'commit', '-S', '-m', 'Create new release'])
    subprocess.check_call(['git', 'push'])


    # Create branch
    if not is_snapshot:
        subprocess.check_call(['git', 'branch', new_version])
        subprocess.check_call(['git', 'push', '-u', 'origin', new_version])
    time.sleep(6)


    # Create release
    response = requests.post(
        'https://api.github.com/repos/%s/%s/releases' % (
            github_owner, REPO_NAME),
        headers={
            'Authorization': 'token %s' % github_token,
            'Content-type': 'application/json',
        },
        data=json.dumps({
            'tag_name': new_version,
            'name': '%s v%s' % (REPO_NAME, new_version),
            'body': release_body,
            'prerelease': is_snapshot,
            'target_commitish': 'master' if is_snapshot else new_version,
        }),
    )

    if response.status_code != 201:
        print('Failed to create release on github')
        print(response.json())
        sys.exit(1)

    subprocess.check_call(['git', 'pull'])
    subprocess.check_call(['git', 'push', '--tags'])
    time.sleep(6)


    # Create gitlab release
    response = requests.post(
        'https://' + gitlab_host + '/api/v4/projects' + \
        '/%s%%2F%s/repository/tags/%s/release' % (
            github_owner, REPO_NAME, new_version),
        headers={
            'Private-Token': gitlab_token,
            'Content-type': 'application/json',
        },
        data=json.dumps({
            'tag_name': new_version,
            'description': release_body,
        }),
    )

    if response.status_code != 201:
        print('Failed to create release on gitlab')
        print(response.json())
        sys.exit(1)


if cmd == 'build' or cmd == 'build-upload':
    is_snapshot = 'snapshot' in cur_version
    pacur_path = TEST_PACUR_PATH if is_snapshot else STABLE_PACUR_PATH


    # Get sha256 sum
    archive_name = '%s.tar.gz' % cur_version
    archive_path = os.path.join(os.path.sep, 'tmp', archive_name)
    if os.path.isfile(archive_path):
        os.remove(archive_path)
    wget('https://github.com/%s/%s/archive/%s' % (
        github_owner, REPO_NAME, archive_name),
        output=archive_name,
        cwd=os.path.join(os.path.sep, 'tmp'),
    )
    archive_sha256_sum = subprocess.check_output(
        ['sha256sum', archive_path]).split()[0]
    os.remove(archive_path)


    # Update sha256 sum and pkgver in PKGBUILD
    for target in BUILD_TARGETS:
        pkgbuild_path = os.path.join(pacur_path, target, 'PKGBUILD')

        with open(pkgbuild_path, 'r') as pkgbuild_file:
            pkgbuild_data = re.sub(
                'pkgver="(.*)"',
                'pkgver="%s"' % cur_version,
                pkgbuild_file.read(),
                count=1,
            )
            pkgbuild_data = re.sub(
                '"[a-f0-9]{64}"',
                '"%s"' % archive_sha256_sum.decode('utf-8'),
                pkgbuild_data,
                count=1,
            )

        with open(pkgbuild_path, 'w') as pkgbuild_file:
            pkgbuild_file.write(pkgbuild_data)


    # Run pacur project build
    for build_target in BUILD_TARGETS:
        subprocess.check_call(
            ['sudo', 'pacur', 'project', 'build', build_target],
            cwd=pacur_path,
        )


if cmd == 'upload' or cmd == 'build-upload':
    is_snapshot = 'snapshot' in cur_version
    pacur_path = TEST_PACUR_PATH if is_snapshot else STABLE_PACUR_PATH


    # Get release id
    release_id = None
    response = requests.get(
        'https://api.github.com/repos/%s/%s/releases' % (
            github_owner, REPO_NAME),
        headers={
            'Authorization': 'token %s' % github_token,
            'Content-type': 'application/json',
        },
    )

    for release in response.json():
        if release['tag_name'] == cur_version:
            release_id = release['id']

    if not release_id:
        print('Version does not exists in github')
        sys.exit(1)


    # Run pacur project build
    subprocess.check_call(
        ['sudo', 'pacur', 'project', 'repo'],
        cwd=pacur_path,
    )


    # Sync mirror
    subprocess.check_call([
        'mc',
        'mirror',
        '--remove',
        '--overwrite',
        '--md5',
        'mirror',
        'repo-east/unstable',
    ], cwd=pacur_path)

    subprocess.check_call([
        'mc',
        'mirror',
        '--remove',
        '--overwrite',
        '--md5',
        'mirror',
        'repo-west/unstable',
    ], cwd=pacur_path)


    # Add to github
    for name, path in iter_packages():
        post_git_asset(release_id, name, path)


if cmd == 'upload-github':
    is_snapshot = 'snapshot' in cur_version


    # Get release id
    release_id = None
    response = requests.get(
        'https://api.github.com/repos/%s/%s/releases' % (
            github_owner, REPO_NAME),
        headers={
            'Authorization': 'token %s' % github_token,
            'Content-type': 'application/json',
        },
    )

    for release in response.json():
        if release['tag_name'] == cur_version:
            release_id = release['id']

    if not release_id:
        print('Version does not exists in github')
        sys.exit(1)


    # Add to github
    for name, path in iter_packages():
        post_git_asset(release_id, name, path)
