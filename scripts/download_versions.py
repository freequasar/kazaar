# TODO: implement as a command in CLI

import requests
import sqlite3
import os.path
import zlib
import io

MANIFESTS = 'https://raw.githubusercontent.com/otomir23/quasar-unlock/refs/heads/main/src/manifests.json'

ROOT_PATH = os.path.dirname(__file__)
DB_PATH = os.path.join(ROOT_PATH, '../kazaar.db')
VERSIONS_PATH = os.path.join(ROOT_PATH, '../versions/')

CHUNK_SIZE = 2 ** 12

INSERT_VERSION_QUERY = '''
INSERT INTO version
(platform, version, crc32, update_script, update_script_signature, update_script_timestamp)
VALUES
(?, ?, ?, ?, ?, ?)
'''


def main() -> None:
    db = sqlite3.connect(DB_PATH)

    with requests.get(MANIFESTS) as r:
        manifests = r.json()

    for platform in manifests.keys():
        manifest = manifests[platform]

        if manifest.get('updateId', '') == '':
            continue

        download_url = manifest['downloadUrl']
        version = manifest['version']
        crc32 = 0

        print(f'[+] Downloading {download_url}')
        with requests.get(download_url, stream=True) as r:
            version_path = os.path.join(VERSIONS_PATH, f'./{version}.zip')

            with open(version_path, 'wb') as f:
                for chunk in r.iter_content(CHUNK_SIZE):
                    f.write(chunk)
                    crc32 = zlib.crc32(chunk, crc32)

        crc32 %= 0xffffffff

        update_script = manifest.get('updateScript')
        update_script_signature = manifest.get('updateScriptSign')
        update_script_timestamp = manifest.get('updateScriptTimestamp')

        print(f'[+] Inserting {version} to DB')
        db.execute(INSERT_VERSION_QUERY, [
            platform, version, crc32,
            update_script, update_script_signature,
            update_script_timestamp,
        ])
        db.commit()

    db.close()


if __name__ == "__main__":
    main()
