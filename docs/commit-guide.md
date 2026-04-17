# 커밋 가이드

## 목적

이 문서는 이 저장소에서 커밋 메시지를 작성할 때 지켜야 하는
최소 규칙을 정리합니다.

## 기본 규칙

- 커밋 메시지는 **Conventional Commits** 형식을 사용합니다.
- 제목은 한 줄로 짧고 명확하게 작성합니다.
- 제목은 명령형 현재 시제를 사용합니다.
- 제목 첫 글자는 불필요하게 대문자로 시작하지 않습니다.
- 제목 끝에 마침표를 붙이지 않습니다.

## 기본 형식

```text
<type>(<scope>): <subject>
```

본문이 필요하면 아래 형식을 사용합니다.

```text
<type>(<scope>): <subject>

<body>
```

## 권장 type

- `feat`: 사용자 기능 추가
- `fix`: 버그 수정
- `refactor`: 동작 변경 없는 구조 개선
- `test`: 테스트 추가 또는 수정
- `docs`: 문서 변경
- `chore`: 빌드, 설정, 의존성, 보조 작업

## scope 기준

scope는 현재 디렉토리 구조와 변경 범위를 기준으로 작성합니다.

- `domain`
- `cmd`
- `docs`

필요하면 실제 변경 패키지나 모듈 이름을 사용해도 됩니다.

## 제목 작성 규칙

- 가능하면 50자 안팎으로 유지합니다.
- 무엇을 바꿨는지 바로 드러나야 합니다.
- 구현 상세보다 변경 의도를 우선합니다.

예시:

```text
refactor(domain): rename PlainEntry to Plain
fix(cmd): return non-zero on invalid input
docs(docs): rename commit guide
```

## 본문 작성 규칙

- 왜 바꿨는지와 무엇이 달라졌는지 적습니다.
- 각 줄은 72자 안팎으로 감쌉니다.
- breaking change나 이슈 번호는 footer로 분리할 수 있습니다.

## 참고 자료

- Conventional Commits: <https://www.conventionalcommits.org/en/about>
- Git `git-commit`: <https://git-scm.com/docs/git-commit>
