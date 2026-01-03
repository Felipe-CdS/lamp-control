# lamp-control

## Procedure based on the following explanation:
> Implementation of the TP-Link Klap Home Protocol.
>
>Encryption/Decryption methods based on the works of
Simon Wilkinson and Chris Weeldon
>
>Klap devices that have never been connected to the kasa
cloud should work with blank credentials.
Devices that have been connected to the kasa cloud will
switch intermittently between the users cloud credentials
and default kasa credentials that are hardcoded.
This appears to be an issue with the devices.
>
>The protocol works by doing a two stage handshake to obtain
and encryption key and session id cookie.
>
>Authentication uses an auth_hash which is
md5(md5(username),md5(password))
>
>handshake1: client sends a random 16 byte local_seed to the
device and receives a random 16 bytes remote_seed, followed
by sha256(local_seed + auth_hash).  It also returns a
TP_SESSIONID in the cookie header.  This implementation
then checks this value against the possible auth_hashes
described above (user cloud, kasa hardcoded, blank).  If it
finds a match it moves onto handshake2
>
>handshake2: client sends sha25(remote_seed + auth_hash) to
the device along with the TP_SESSIONID.  Device responds with
200 if successful.  It generally will be because this
implementation checks the auth_hash it received during handshake1
>
>encryption: local_seed, remote_seed and auth_hash are now used
for encryption.  The last 4 bytes of the initialization vector
are used as a sequence number that increments every time the
client calls encrypt and this sequence number is sent as a
url parameter to the device along with the encrypted payload
>
>https://gist.github.com/chriswheeldon/3b17d974db3817613c69191c0480fe55
> 
>https://github.com/python-kasa/python-kasa/pull/117



## 1) Handshake 1

### One-line testing cmd:
``` bash
dd if=/dev/urandom bs=16 count=1 status=none | curl -X POST --data-binary @- -i --output "@output.bin" http://192.168.1.2:80/app/handshake1
```
### Step by step procedure:
- Create a random 16 byte token (local_seed):
- ```
  dd if=/dev/urandom of=handshake1_payload.bin bs=16 count=1
  ```
- Send a POST with the token:
- ```
  curl --request POST --data-binary "@handshake1_payload.bin" -H "Content-Type: application/octet-stream" -i http://192.168.1.2:80/app/handshake1
  ```
- The response should be something like:
- ```
  HTTP/1.1 200 OK
  Set-Cookie: TP_SESSIONID=295D13933FC214B1A0C996B22B0686A4;TIMEOUT=86400
  Server: SHIP 2.0
  Content-Length: 48
  Content-Type: text/html
  
  dc��ĕ��6�ԥ�J���S'W�Z�Q�a�X	�F�ow�v��L4%
  ```
- Response will return some bytes on the body.
  - `local_seed := handshake1_payload.bin // The 16-byte token you sent on the request`
  - `remote_seed := response[0:16]`
  - `server_hash := response[16:]`

- Generate an `auth_hash` based on the TP_link Cloud credentials:
  - `md5(md5(username) + md5(password))`

- Generate a `local_seed_auth_hash`:
  - `sha256(local_seed + remote_seed + auth_hash)`

- If `local_seed_auth_hash == server_hash` then handshake1 worked.

## 2) Handshake 2
