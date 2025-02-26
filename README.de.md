# gh-otui

/oˈtuː.i/ wird gelesen.

gh-otui = gh + org + tui

![otui](https://github.com/user-attachments/assets/0c7626eb-c639-4f4c-86e1-b4ba6dab5bec)  
(Das Repository, das im GIF angezeigt wird, ist öffentlich und gehört zu der Organisation, zu der ich gehöre.)

gh-otui ist ein CLI-Tool, das gh mit ghq und einem fuzzy finder (peco, fzf) kombiniert.  
Es ermöglicht das Durchsuchen und Anzeigen von Organisationen und eigenen Repositories unter Verwendung der fuzzy finder Logik und das Klonen mit ghq. Besonders wenn Sie an mehreren Repositories gleichzeitig arbeiten, ist es praktisch, dass Sie nur den Repository-Namen kennen müssen, um das Klonen ausschließlich über die CLI abzuschließen.

## Funktionen

- Anzeige einer Liste der GitHub-Organisationen und der eigenen Repositories
- Interaktive Repository-Auswahl mit fuzzy finder
- Klonen des ausgewählten Repositories mit ghq (für noch nicht geklonte Repositories)
- Visuelle Anzeige geklonter Repositories (✓-Markierung)

## Vorrausgesetzte Tools

- [GitHub CLI](https://cli.github.com/) (gh)
- [ghq](https://github.com/x-motemen/ghq)
- [peco](https://github.com/peco/peco)
  - oder [fzf](https://github.com/junegunn/fzf). Durch Setzen der Umgebungsvariablen `GH_OTUI_SELECTOR` auf `fzf` können Sie fzf verwenden. Wenn keine Umgebungsvariable angegeben ist, wird das installiert, was sowohl bei peco als auch bei fzf vorhanden ist. Wenn beide installiert sind, hat peco Vorrang.
  
## Installation

```bash
gh extension install n3xem/gh-otui
```

## Verwendung

1. Erstellen Sie einen Cache für die Repositories der Organisation, der Sie angehören:

```bash
gh otui --cache
```

Der Cache wird in `~/.config/gh/extensions/gh-otui/cache.json` gespeichert.

2. Führen Sie den folgenden Befehl aus:

```bash
gh otui
```

3. Wählen Sie im fuzzy finder Interface das gewünschte Repository aus
   - Die ✓-Markierung zeigt Repositories an, die bereits geklont wurden
   - Wenn Sie ein noch nicht geklontes Repository auswählen, wird es mit ghq geklont
   - Ob ein Repository geklont ist, wird durch Überprüfung des Pfades von `ghq root` festgestellt

4. Der lokale Pfad des ausgewählten Repositories wird auf der Standardausgabe ausgegeben.
   - Es ist praktisch, dies zusammen mit dem cd-Befehl zu verwenden, um schnell dorthin zu navigieren.
   - Beispiel: `cd $(gh otui)`

## Ausgabeformat

Die Repositories werden im folgenden Format angezeigt:

- ✓: Markierung für geklonte Repositories
- organization-name: Name der GitHub-Organisation
- repository-name: Name des Repositories