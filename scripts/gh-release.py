import os
import sys

import requests

from requests.auth import HTTPBasicAuth


OWNER = 'mzbaulhaque'
REPO = 'gomage'
AUTH = HTTPBasicAuth(OWNER, os.environ['GH_ACCESS_TOKEN'])
HEADERS = {
    'Accept': 'application/vnd.github.v3+json',
}
PLATFORMS = [
    ('darwin', 'amd64'),
    ('freebsd', 'amd64'),
    ('linux', 'amd64'),
    ('linux', 'arm64'),
    ('windows', 'amd64'),
]

def create_release(tag, name, notes):
    body = {
        'tag_name': tag,
        'name': name,
        'body': notes,
    }
    res = requests.post(f'https://api.github.com/repos/{OWNER}/{REPO}/releases', auth=AUTH, json=body, headers=headers)
    res.raise_for_status()
    res = res.json()

    return res['id']

def upload_assets(release_id, tag, version):
    headers = dict(**HEADERS, **{
        'Content-Type': 'application/zip',
    })
    upload_url = f'https://uploads.github.com/repos/{OWNER}/{REPO}/releases/{release_id}/assets'
    for p in PLATFORMS:
        filename = f'gomage-{tag}.{p[0]}-{p[1]}.tar.gz'
        with open(f'dist/{filename}', 'rb') as f:
            res = requests.post(f'{upload_url}?name={filename}', auth=AUTH, headers=headers, data=f)

def get_release_info(tag, version):
    info = {'name': '', 'notes': ''}
    with open('CHANGELOG.md', 'w') as f:
        cl_content = f.readlines()
        notes = ''
        start_note = False

        for line in cl_content:
            if line.startswith(f'## [{version}]'):
                info['name'] = f'tag - {line.split(" ")[-1]}'
                start_note = True
            elif line.startswith('## [') and start_note:
                break
            else:
                notes += line

        info['notes'] = notes.strip()

def main():
    args = sys.argv[1:]
    tag = args[1]
    version = tag.split('v')[1]
    release_info = get_release_info(tag, version)
    release_id = create_release(tag, release_info['name'], release_info['notes'])
    upload_assets(release_id)

    return 0

if __name__ == '__main__':
    sys.exit(main())
