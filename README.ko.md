# gh-otui

/oˈtuː.i/ 라고 읽습니다.

gh-otui = gh + org + tui

![otui](https://github.com/user-attachments/assets/0c7626eb-c639-4f4c-86e1-b4ba6dab5bec)
(화면에 표시된 GIF 리포지토리는 모두 제가 소속된 조직의 공개된 것입니다.)

gh-otui는 gh와 ghq, 퍼지 파인더(peco, fzf)를 조합한 CLI 도구입니다.  
조직이나 자신의 리포지토리를 퍼지 파인더의 구조를 사용하여 가로질러 검색하고 열람하며, ghq를 통해 클론할 수 있습니다. 특히 여러 리포지토리를 가로질러 개발하는 경우, 리포지토리 이름만 알고 있다면 CLI만으로 클론을 완료할 수 있어서 매우 편리합니다.

## 기능

- GitHub의 조직 및 자신의 리포지토리 목록 표시
- 퍼지 파인더를 사용한 대화형 리포지토리 선택
- 선택한 리포지토리의 ghq에 의한 클론(미클론의 경우)
- 클론된 리포지토리의 시각적 표시(✓ 마크)

## 전제 도구

- [GitHub CLI](https://cli.github.com/) (gh)
- [ghq](https://github.com/x-motemen/ghq)
- [peco](https://github.com/peco/peco)
  - 또는 [fzf](https://github.com/junegunn/fzf). 환경 변수 `GH_OTUI_SELECTOR` 를 `fzf`로 설정하면 fzf를 사용할 수 있습니다. 환경 변수가 지정되지 않은 경우, peco와 fzf 중 설치된 것을 사용합니다. 두 개가 모두 설치되어 있으면 peco가 우선됩니다.

## 설치

```bash
gh extension install n3xem/gh-otui
```

## 사용 방법

1. 소속한 조직의 리포지토리 캐시를 생성합니다:

```bash
gh otui --cache
```

캐시는 `~/.config/gh/extensions/gh-otui/cache.json` 에 저장됩니다.

2. 다음 명령어를 실행합니다:

```bash
gh otui
```

3. 퍼지 파인더 인터페이스에서 원하는 리포지토리를 선택합니다.
   - ✓ 마크는 이미 클론된 리포지토리를 나타냅니다.
   - 미클론의 리포지토리를 선택하면 ghq에 의해 클론이 진행됩니다.
   - 클론 여부는 `ghq root`의 경로를 확인하여 결정됩니다.

4. 선택한 리포지토리의 로컬 경로가 표준 출력됩니다.
   - cd 명령어와 연계하여 사용하면 바로 이동할 수 있어 편리합니다.
   - 예: `cd $(gh otui)`

## 출력 형식

리포지토리는 다음 형식으로 표시됩니다:

- ✓: 클론된 리포지토리를 나타내는 마크
- organization-name: GitHub의 조직명
- repository-name: 리포지토리명