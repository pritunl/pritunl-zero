import datetime
import re
import sys
import subprocess
import math
import json
import os

CONSTANTS_PATH = 'constants/constants.go'
STABLE_PACUR_PATH = '../pritunl-pacur'
TEST_PACUR_PATH = '../pritunl-pacur-test'
BUILD_KEYS_PATH = 'build_keys.json'
BUILD_TARGETS = ('pritunl-zero',)

cur_date = datetime.datetime.utcnow()

with open(BUILD_KEYS_PATH, 'r') as build_keys_file:
    build_keys = json.loads(build_keys_file.read().strip())
    mirror_url = build_keys['mirror_url']
    test_mirror_url = build_keys['test_mirror_url']

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

cmd = sys.argv[1]

with open(CONSTANTS_PATH, 'r') as constants_file:
    cur_version = re.findall('= "(.*?)"', constants_file.read())[0]

if cmd == 'version':
    print get_ver(sys.argv[2])

elif cmd == 'set-version':
    new_version = get_ver(sys.argv[2])

    with open(CONSTANTS_PATH, 'r') as constants_file:
        constants_data = constants_file.read()

    with open(CONSTANTS_PATH, 'w') as constants_file:
        constants_file.write(re.sub(
            '(= ".*?")',
            '= "%s"' % new_version,
            constants_data,
            count=1,
        ))

    subprocess.check_call(['git', 'reset', 'HEAD', '.'])
    subprocess.check_call(['git', 'add', CONSTANTS_PATH])
    subprocess.check_call(['git', 'commit', '-S', '-m', 'Create new release'])
    subprocess.check_call(['git', 'push'])

elif cmd == 'build' or cmd == 'build-test':
    if cmd == 'build':
        pacur_path = STABLE_PACUR_PATH
    else:
        pacur_path = TEST_PACUR_PATH

    for target in BUILD_TARGETS:
        pkgbuild_path = os.path.join(pacur_path, target, 'PKGBUILD')

        with open(pkgbuild_path, 'r') as pkgbuild_file:
            pkgbuild_data = re.sub(
                'pkgver="(.*)"',
                'pkgver="%s"' % cur_version,
                pkgbuild_file.read(),
            )

        with open(pkgbuild_path, 'w') as pkgbuild_file:
            pkgbuild_file.write(pkgbuild_data)

    for build_target in BUILD_TARGETS:
        subprocess.check_call(
            ['sudo', 'pacur', 'project', 'build', build_target],
            cwd=pacur_path,
        )

elif cmd == 'upload' or cmd == 'upload-test':
    if cmd == 'upload':
        mirror_urls = mirror_url
        pacur_path = STABLE_PACUR_PATH
    else:
        mirror_urls = test_mirror_url
        pacur_path = TEST_PACUR_PATH

    subprocess.check_call(
        ['sudo', 'pacur', 'project', 'repo'],
        cwd=pacur_path,
    )

    for mir_url in mirror_urls:
        subprocess.check_call([
            'rsync',
            '--human-readable',
            '--archive',
            '--progress',
            '--delete',
            '--acls',
            'mirror/',
            mir_url,
        ], cwd=pacur_path)
