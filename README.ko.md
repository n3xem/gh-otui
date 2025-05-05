# gh-otui

/oˈtuː.i/ 라고 읽습니다.

gh-otui = gh + org + tui

![otui](https://github.com/user-attachments/assets/0c7626eb-c639-4f4c-86e1-b4ba6dab5bec)  
(애니메이션으로 표시된 리포지토리는 모두 제가 속한 조직의 공개 리포지토리입니다)

gh-otui는 gh와 ghq, 퍼지 파인더(peco, fzf)를 결합한 CLI 도구입니다.  
조직이나 자신의 리포지토리를 퍼지 파인더의 방식을 사용하여 가로질러 검색 및 열람하며, ghq로 클론할 수 있습니다. 특히 여러 리포지토리를 가로질러 개발하는 경우, 리포지토리 이름만 알고 있으면 CLI만으로 클론을 완료할 수 있어 편리합니다.

## 기능

- GitHub의 조직, 자신의 리포지토리 목록 표시
- 퍼지 파인더를 사용한 대화형 리포지토리 선택
- 선택한 리포지토리의 ghq에 의한 클론(미클론 경우)
- 클론된 리포지토리의 시각적 표시(✓ 마크)

## 전제 도구

- [GitHub CLI](https://cli.github.com/) (gh)
- [ghq](https://github.com/x-motemen/ghq)
- [peco](https://github.com/peco/peco)
  - 또는 [fzf](https://github.com/junegunn/fzf). 환경변수 `GH_OTUI_SELECTOR`를 `fzf`로 설정하면 fzf를 사용할 수 있습니다. 환경변수 지정이 없는 경우, peco와 fzf가 설치된 경우에는 설치된 도구 중 peco가 우선적으로 사용됩니다.
  
## 설치

```bash
gh extension install n3xem/gh-otui
```

## 사용 방법

1. gh otui 명령어를 실행하기만 하면 됩니다. 처음 실행 시 가져올 리포지토리 목록을 저장한 캐시가 생성됩니다.

```bash
gh otui
```

2. 퍼지 파인더 인터페이스에서 원하는 리포지토리를 선택합니다.
   - ✓ 마크는 이미 클론된 리포지토리를 나타냅니다.
   - 미클론 리포지토리를 선택하면 ghq로 클론됩니다.
   - 클론 여부는 `ghq root`의 경로를 확인하여 판단됩니다.

3. 선택한 리포지토리의 로컬 경로가 표준 출력으로 표시됩니다.
   - cd 명령어와 연계하여 사용하면 즉시 이동이 가능하여 편리합니다.
   - 예: `cd $(gh otui)`

## 출력 형식

리포지토리는 다음 형식으로 표시됩니다:

- ✓: 클론된 리포지토리를 나타내는 마크
- organization-name: GitHub 조직 이름
- repository-name: 리포지토리 이름

## 캐시에 대하여

gh-otui는 다음과 같은 캐시 구조를 사용합니다:

- **캐시 저장 위치**: `~/.config/gh/extensions/gh-otui/`
- **유효 기간**: 1시간(최종 업데이트로부터 1시간이 경과하면 오래된 것으로 판단되어 백그라운드에서 자동으로 업데이트됩니다)
- **메타데이터 파일**: `_md.json` - 캐시의 최종 업데이트 시간을 저장
- **호스트 디렉토리**: 각 GitHub 호스트(예: `github.com`)마다 디렉토리 생성
- **조직 파일**: 각 조직마다 `{organization}.json` 파일을 생성, 리포지토리 정보를 저장

캐시 업데이트는 다음과 같은 경우에 진행됩니다:
1. 최초 실행 시(캐시가 존재하지 않는 경우)
2. 캐시 유효 기간(1시간)이 만료된 경우(백그라운드에서 자동 업데이트)

캐시 삭제: `gh otui clear` 명령어로 캐시 디렉토리를 삭제할 수 있습니다.