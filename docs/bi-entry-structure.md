# BI Entry Structure

BI에서는 버저닝 오브젝트와 비버저닝 오브젝트(`Unversioned Object`)의 형태가 다르다.

## 비버저닝 오브젝트

비버저닝 오브젝트는 BI에 `plain entry`만 존재한다.

`radosgw-admin bi list --bucket=test --object=b.txt`로 조회하면, 비버저닝 오브젝트는 아래처럼 `plain` 타입 엔트리 하나만 보인다.

```json
[
    {
        "type": "plain",
        "idx": "b.txt",
        "entry": {
            "name": "b.txt",
            "instance": "",
            "ver": {
                "pool": 186,
                "epoch": 1338
            },
            "locator": "",
            "exists": true,
            "meta": {
                "category": 1,
                "size": 4,
                "mtime": "2026-03-06T03:34:11.918188Z",
                "etag": "3a029f04d76d32e79367c4b3255dda4d",
                "storage_class": "",
                "owner": "test",
                "owner_display_name": "test",
                "content_type": "text/plain",
                "accounted_size": 4,
                "user_data": "",
                "appendable": false
            },
            "tag": "c0688820-3bc9-43a2-8eb2-e688ea8bd667.1324011.5028615704413862483",
            "flags": 0,
            "pending_map": [],
            "versioned_epoch": 0
        }
    }
]
```

## 버저닝 오브젝트

버저닝 오브젝트는 BI에 아래 엔트리들이 존재한다.

- `1 head plain entry`
- `1 olh entry`
- 정상적인 경우 각 version마다 `plain entry + instance entry` 한 쌍

이 규칙과 다르게 보이는 경우는 버그 상황으로 간주한다.

여기서 `head plain entry`는 현재 head 상태를 나타내는 placeholder 성격의 엔트리다.

`olh`는 Object Logical Head의 약자로, 최신 instance를 가리키는 용도로 사용된다. 최신 상태가 delete marker인 경우에도 해당 delete-marker instance를 가리킨다.

`radosgw-admin bi list --bucket=test --object=b.txt`로 조회하면, 버저닝 오브젝트는 `head plain entry`, 버전별 `plain/instance entry`, 그리고 `olh entry`를 함께 가진다.

```json
[
    {
        "type": "plain",
        "idx": "b.txt",
        "entry": {
            "name": "b.txt",
            "instance": "",
            "ver": {
                "pool": -1,
                "epoch": 0
            },
            "locator": "",
            "exists": false,
            "meta": {
                "category": 0,
                "size": 0,
                "mtime": "0.000000",
                "etag": "",
                "storage_class": "",
                "owner": "",
                "owner_display_name": "",
                "content_type": "",
                "accounted_size": 0,
                "user_data": "",
                "appendable": false
            },
            "tag": "",
            "flags": 8,
            "pending_map": [],
            "versioned_epoch": 0
        }
    },
    {
        "type": "plain",
        "idx": "b.txt\u0000v913\u0000iGAhxeKhXP.WC7-4eNqJdvgAQV-CGCFq",
        "entry": {
            "name": "b.txt",
            "instance": "GAhxeKhXP.WC7-4eNqJdvgAQV-CGCFq",
            "ver": {
                "pool": 186,
                "epoch": 1147
            },
            "locator": "",
            "exists": true,
            "meta": {
                "category": 1,
                "size": 4,
                "mtime": "2026-03-06T01:20:38.700296Z",
                "etag": "3a029f04d76d32e79367c4b3255dda4d",
                "storage_class": "",
                "owner": "test",
                "owner_display_name": "test",
                "content_type": "text/plain",
                "accounted_size": 4,
                "user_data": "",
                "appendable": false
            },
            "tag": "c0688820-3bc9-43a2-8eb2-e688ea8bd667.1324011.12576304871632696882",
            "flags": 3,
            "pending_map": [],
            "versioned_epoch": 2
        }
    },
    {
        "type": "instance",
        "idx": "�1000_b.txt\u0000iGAhxeKhXP.WC7-4eNqJdvgAQV-CGCFq",
        "entry": {
            "name": "b.txt",
            "instance": "GAhxeKhXP.WC7-4eNqJdvgAQV-CGCFq",
            "ver": {
                "pool": 186,
                "epoch": 1147
            },
            "locator": "",
            "exists": true,
            "meta": {
                "category": 1,
                "size": 4,
                "mtime": "2026-03-06T01:20:38.700296Z",
                "etag": "3a029f04d76d32e79367c4b3255dda4d",
                "storage_class": "",
                "owner": "test",
                "owner_display_name": "test",
                "content_type": "text/plain",
                "accounted_size": 4,
                "user_data": "",
                "appendable": false
            },
            "tag": "c0688820-3bc9-43a2-8eb2-e688ea8bd667.1324011.12576304871632696882",
            "flags": 3,
            "pending_map": [],
            "versioned_epoch": 2
        }
    },
    {
        "type": "olh",
        "idx": "�1001_b.txt",
        "entry": {
            "key": {
                "name": "b.txt",
                "instance": "GAhxeKhXP.WC7-4eNqJdvgAQV-CGCFq"
            },
            "delete_marker": false,
            "epoch": 2,
            "pending_log": [],
            "tag": "h5hx9nksbsnc09zyqvlxv0gnkmihqocr",
            "exists": true,
            "pending_removal": false
        }
    }
]
```

버저닝 오브젝트를 1회 업로드한 뒤 삭제하면, delete-marker version이 추가되고 `olh`는 그 delete-marker instance를 가리킨다.

```json
[
    {
        "type": "plain",
        "idx": "b.txt",
        "entry": {
            "name": "b.txt",
            "instance": "",
            "ver": {
                "pool": -1,
                "epoch": 0
            },
            "locator": "",
            "exists": false,
            "meta": {
                "category": 0,
                "size": 0,
                "mtime": "0.000000",
                "etag": "",
                "storage_class": "",
                "owner": "",
                "owner_display_name": "",
                "content_type": "",
                "accounted_size": 0,
                "user_data": "",
                "appendable": false
            },
            "tag": "",
            "flags": 8,
            "pending_map": [],
            "versioned_epoch": 0
        }
    },
    {
        "type": "plain",
        "idx": "b.txt\u0000v912\u0000i3BC0g1tGmmeLQnurWUtXeRy8gX0DMHE",
        "entry": {
            "name": "b.txt",
            "instance": "3BC0g1tGmmeLQnurWUtXeRy8gX0DMHE",
            "ver": {
                "pool": -1,
                "epoch": 0
            },
            "locator": "",
            "exists": false,
            "meta": {
                "category": 0,
                "size": 0,
                "mtime": "2026-03-06T04:11:07.657765Z",
                "etag": "",
                "storage_class": "",
                "owner": "test",
                "owner_display_name": "test",
                "content_type": "",
                "accounted_size": 0,
                "user_data": "",
                "appendable": false
            },
            "tag": "delete-marker",
            "flags": 7,
            "pending_map": [],
            "versioned_epoch": 3
        }
    },
    {
        "type": "plain",
        "idx": "b.txt\u0000v913\u0000iCTAKqJ1kcR500Nq435affizMW8bu0P8",
        "entry": {
            "name": "b.txt",
            "instance": "CTAKqJ1kcR500Nq435affizMW8bu0P8",
            "ver": {
                "pool": 186,
                "epoch": 1388
            },
            "locator": "",
            "exists": true,
            "meta": {
                "category": 1,
                "size": 4,
                "mtime": "2026-03-06T04:11:01.457183Z",
                "etag": "3a029f04d76d32e79367c4b3255dda4d",
                "storage_class": "",
                "owner": "test",
                "owner_display_name": "test",
                "content_type": "text/plain",
                "accounted_size": 4,
                "user_data": "",
                "appendable": false
            },
            "tag": "c0688820-3bc9-43a2-8eb2-e688ea8bd667.1324011.18045272437437620319",
            "flags": 1,
            "pending_map": [],
            "versioned_epoch": 2
        }
    },
    {
        "type": "instance",
        "idx": "�1000_b.txt\u0000i3BC0g1tGmmeLQnurWUtXeRy8gX0DMHE",
        "entry": {
            "name": "b.txt",
            "instance": "3BC0g1tGmmeLQnurWUtXeRy8gX0DMHE",
            "ver": {
                "pool": -1,
                "epoch": 0
            },
            "locator": "",
            "exists": false,
            "meta": {
                "category": 0,
                "size": 0,
                "mtime": "2026-03-06T04:11:07.657765Z",
                "etag": "",
                "storage_class": "",
                "owner": "test",
                "owner_display_name": "test",
                "content_type": "",
                "accounted_size": 0,
                "user_data": "",
                "appendable": false
            },
            "tag": "delete-marker",
            "flags": 7,
            "pending_map": [],
            "versioned_epoch": 3
        }
    },
    {
        "type": "instance",
        "idx": "�1000_b.txt\u0000iCTAKqJ1kcR500Nq435affizMW8bu0P8",
        "entry": {
            "name": "b.txt",
            "instance": "CTAKqJ1kcR500Nq435affizMW8bu0P8",
            "ver": {
                "pool": 186,
                "epoch": 1388
            },
            "locator": "",
            "exists": true,
            "meta": {
                "category": 1,
                "size": 4,
                "mtime": "2026-03-06T04:11:01.457183Z",
                "etag": "3a029f04d76d32e79367c4b3255dda4d",
                "storage_class": "",
                "owner": "test",
                "owner_display_name": "test",
                "content_type": "text/plain",
                "accounted_size": 4,
                "user_data": "",
                "appendable": false
            },
            "tag": "c0688820-3bc9-43a2-8eb2-e688ea8bd667.1324011.18045272437437620319",
            "flags": 1,
            "pending_map": [],
            "versioned_epoch": 2
        }
    },
    {
        "type": "olh",
        "idx": "�1001_b.txt",
        "entry": {
            "key": {
                "name": "b.txt",
                "instance": "3BC0g1tGmmeLQnurWUtXeRy8gX0DMHE"
            },
            "delete_marker": true,
            "epoch": 3,
            "pending_log": [],
            "tag": "zeaazy8aey9eejp58d0hfgo3k4utj37s",
            "exists": true,
            "pending_removal": false
        }
    }
]
```

버전이 2개인 경우에는 version별 `plain/instance entry`가 2쌍 생기고, `olh`는 최신 version의 instance를 가리킨다.

```json
[
    {
        "type": "plain",
        "idx": "b.txt",
        "entry": {
            "name": "b.txt",
            "instance": "",
            "ver": {
                "pool": -1,
                "epoch": 0
            },
            "locator": "",
            "exists": false,
            "meta": {
                "category": 0,
                "size": 0,
                "mtime": "0.000000",
                "etag": "",
                "storage_class": "",
                "owner": "",
                "owner_display_name": "",
                "content_type": "",
                "accounted_size": 0,
                "user_data": "",
                "appendable": false
            },
            "tag": "",
            "flags": 8,
            "pending_map": [],
            "versioned_epoch": 0
        }
    },
    {
        "type": "plain",
        "idx": "b.txt\u0000v912\u0000ik8xZlwCiGIG7flEibzi1UN0EpiuJH-E",
        "entry": {
            "name": "b.txt",
            "instance": "k8xZlwCiGIG7flEibzi1UN0EpiuJH-E",
            "ver": {
                "pool": 186,
                "epoch": 1121
            },
            "locator": "",
            "exists": true,
            "meta": {
                "category": 1,
                "size": 4,
                "mtime": "2026-03-06T02:10:28.562296Z",
                "etag": "3a029f04d76d32e79367c4b3255dda4d",
                "storage_class": "",
                "owner": "test",
                "owner_display_name": "test",
                "content_type": "text/plain",
                "accounted_size": 4,
                "user_data": "",
                "appendable": false
            },
            "tag": "c0688820-3bc9-43a2-8eb2-e688ea8bd667.1324011.16080191998712481675",
            "flags": 3,
            "pending_map": [],
            "versioned_epoch": 3
        }
    },
    {
        "type": "plain",
        "idx": "b.txt\u0000v913\u0000iGAhxeKhXP.WC7-4eNqJdvgAQV-CGCFq",
        "entry": {
            "name": "b.txt",
            "instance": "GAhxeKhXP.WC7-4eNqJdvgAQV-CGCFq",
            "ver": {
                "pool": 186,
                "epoch": 1147
            },
            "locator": "",
            "exists": true,
            "meta": {
                "category": 1,
                "size": 4,
                "mtime": "2026-03-06T01:20:38.700296Z",
                "etag": "3a029f04d76d32e79367c4b3255dda4d",
                "storage_class": "",
                "owner": "test",
                "owner_display_name": "test",
                "content_type": "text/plain",
                "accounted_size": 4,
                "user_data": "",
                "appendable": false
            },
            "tag": "c0688820-3bc9-43a2-8eb2-e688ea8bd667.1324011.12576304871632696882",
            "flags": 1,
            "pending_map": [],
            "versioned_epoch": 2
        }
    },
    {
        "type": "instance",
        "idx": "�1000_b.txt\u0000iGAhxeKhXP.WC7-4eNqJdvgAQV-CGCFq",
        "entry": {
            "name": "b.txt",
            "instance": "GAhxeKhXP.WC7-4eNqJdvgAQV-CGCFq",
            "ver": {
                "pool": 186,
                "epoch": 1147
            },
            "locator": "",
            "exists": true,
            "meta": {
                "category": 1,
                "size": 4,
                "mtime": "2026-03-06T01:20:38.700296Z",
                "etag": "3a029f04d76d32e79367c4b3255dda4d",
                "storage_class": "",
                "owner": "test",
                "owner_display_name": "test",
                "content_type": "text/plain",
                "accounted_size": 4,
                "user_data": "",
                "appendable": false
            },
            "tag": "c0688820-3bc9-43a2-8eb2-e688ea8bd667.1324011.12576304871632696882",
            "flags": 1,
            "pending_map": [],
            "versioned_epoch": 2
        }
    },
    {
        "type": "instance",
        "idx": "�1000_b.txt\u0000ik8xZlwCiGIG7flEibzi1UN0EpiuJH-E",
        "entry": {
            "name": "b.txt",
            "instance": "k8xZlwCiGIG7flEibzi1UN0EpiuJH-E",
            "ver": {
                "pool": 186,
                "epoch": 1121
            },
            "locator": "",
            "exists": true,
            "meta": {
                "category": 1,
                "size": 4,
                "mtime": "2026-03-06T02:10:28.562296Z",
                "etag": "3a029f04d76d32e79367c4b3255dda4d",
                "storage_class": "",
                "owner": "test",
                "owner_display_name": "test",
                "content_type": "text/plain",
                "accounted_size": 4,
                "user_data": "",
                "appendable": false
            },
            "tag": "c0688820-3bc9-43a2-8eb2-e688ea8bd667.1324011.16080191998712481675",
            "flags": 3,
            "pending_map": [],
            "versioned_epoch": 3
        }
    },
    {
        "type": "olh",
        "idx": "�1001_b.txt",
        "entry": {
            "key": {
                "name": "b.txt",
                "instance": "k8xZlwCiGIG7flEibzi1UN0EpiuJH-E"
            },
            "delete_marker": false,
            "epoch": 3,
            "pending_log": [],
            "tag": "h5hx9nksbsnc09zyqvlxv0gnkmihqocr",
            "exists": true,
            "pending_removal": false
        }
    }
]
```
