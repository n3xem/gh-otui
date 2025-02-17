# gh-otui

/oˈtuː.i/로 읽힙니다.

gh-otui = gh + org + tui

![otui](https://github.com/user-attachments/assets/0c7626eb-c639-4f4c-86e1-b4ba6dab5bec)
(이 GIF로 표시된 리포지토리는 모두 제가 소속된 조직의 공개 리포지토리입니다.)

gh-otui는 gh와 ghq, fuzz finder(peco, fzf)를 결합한 CLI 도구입니다.  
조직의 리포지토리를 fuzz finder의 구조를 사용하여 가로질러 검색하고 탐색할 수 있으며, ghq를 사용하여 클론할 수 있습니다. 특히 여러 리포지토리를 가로질러 개발하는 경우, 리포지토리 이름만 알면 CLI만으로 클론을 완료할 수 있어 매우 편리합니다.

## 기능

- GitHub의 조직 리포지토리 목록 표시
- fuzz finder를 사용한 상호작용적인 리포지토리 선택
- 선택한 리포지토리의 ghq에 의한 클론(클론되지 않은 경우)
- 클론된 리포지토리의 시각적 표시(✓ 마크)

## 전제 도구

- [GitHub CLI](https://cli.github.com/) (gh)
- [ghq](https://github.com/x-motemen/ghq)
- [peco](https://github.com/peco/peco)
  - 또는 [fzf](https://github.com/junegunn/fzf). 환경 변수 `GH_OTUI_SELECTOR`를 `fzf`로 설정하면 fzf를 사용할 수 있습니다. 환경 변수 지정이 없는 경우 peco와 fzf가 설치된 쪽을 사용합니다. 두 개 모두 설치되어 있는 경우 peco가 우선됩니다.
  
## 설치

```bash
gh extension install n3xem/gh-otui
```

## 사용법

1. 소속된 organization의 리포지토리 캐시를 만듭니다:

```bash
gh otui --cache
```

캐시는 `~/.config/gh/extensions/gh-otui/cache.json`에 저장됩니다.

2. 다음 명령어를 실행합니다:

```bash
gh otui
```

3. fuzz finder 인터페이스에서 원하는 리포지토리를 선택합니다
   - ✓ 마크는 이미 클론된 리포지토리를 나타냅니다
   - 클론되지 않은 리포지토리를 선택하면 ghq에 의해 클론이 진행됩니다
   - 클론 여부는 `ghq root`의 경로를 확인하여 판단됩니다

4. 선택한 리포지토리의 로컬 경로가 표준 출력됩니다.
   - cd 명령어와 연계하여 사용하면 즉시 이동할 수 있어 편리합니다.
   - 예: `cd $(gh otui)`

## 출력 형식

리포지토리는 다음 형식으로 표시됩니다:

- ✓: 클론된 리포지토리를 나타내는 마크
- organization-name: GitHub의 조직 이름
- repository-name: 리포지토리 이름