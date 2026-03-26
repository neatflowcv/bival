# ADR 0003: 버저닝되지 않은 오브젝트의 용어는 `Unversioned Object`로 통일한다

- 상태: 승인됨
- 날짜: 2026-03-26

## 배경

이 프로젝트는 BI 엔트리 그룹을 해석하면서 버저닝 오브젝트와 버저닝되지 않은 오브젝트를 구분한다.

이때 영어 용어 후보로 아래 표현들이 함께 거론됐다.

- `Unversioned Object`
- `NonVersioned Object`
- `Normal Object`

용어를 새로 만들기보다, 실제 Ceph 문맥에서 더 자연스럽고 충돌이 적은 표현을 기준으로 정리할 필요가 있었다.

## 결정

버저닝되지 않은 오브젝트의 프로젝트 표준 용어는 `Unversioned Object`로 한다.

`NonVersioned Object`, `Normal Object`는 표준 용어로 사용하지 않는다.

## 근거

- 현재 프로젝트 코드도 이미 `unversioned_object`를 사용한다.
- Ceph 문서에서는 버저닝이 없는 경우를 `un-versioned`라고 직접 표현한다.
  - `rgw-restore-bucket-index` 문서는 `regular (i.e., un-versioned) buckets`라고 설명한다.
- Ceph의 RGW 버저닝 관련 릴리스 노트에서는 `plain object`라는 표현도
  보이지만, 이는 특정 동작 설명에 가까우며 프로젝트 전반의 분류명으로
  쓰기에는 의미가 좁다.
  - 예: `convert plain object to versioned (with null version) when removing`
- `Normal Object`는 Ceph에서 이미 다른 의미로 사용된다.
  - S3 Object Operations 문서에서 `Appendable Object`의 반대 개념으로 `Normal Object`를 사용한다.
  - 따라서 버저닝 여부를 나타내는 용어로 가져오면 의미 충돌이 생긴다.
- `NonVersioned Object`는 이해는 가능하지만, 확인한 Ceph 문맥에서는 `un-versioned` 쪽이 더 직접적으로 나타난다.

## 결과

긍정적 영향은 아래와 같다.

- 프로젝트 내부 분류명과 현재 코드가 일치한다.
- Ceph 문맥과 연결되는 용어를 사용하면서도 의미 충돌을 줄일 수 있다.
- `plain entry`와 `Normal Object` 같은 다른 기존 용어와 구분이 명확해진다.

주의사항은 아래와 같다.

- Ceph 자료에서 `plain object`라는 표현이 나오더라도, 이 프로젝트의 표준 분류명으로 자동 치환하지는 않는다.
- 한국어 문서에서는 필요에 따라 `비버저닝 오브젝트`라고 쓸 수 있지만, 대응 영어 용어는 `Unversioned Object`로 본다.

## 참고

- Ceph 문서: `rgw-restore-bucket-index`의 `regular (i.e., un-versioned) buckets`
- Ceph 릴리스 노트:
  `rgw: convert plain object to versioned (with null version) when removing`
- Ceph S3 Object Operations: `Appendable Object`와 대비되는 `Normal Object`
