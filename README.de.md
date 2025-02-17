# gh-otui

/oˈtuː.i/ wird so ausgesprochen.

gh-otui = gh + org + tui

![otui](https://github.com/user-attachments/assets/0c7626eb-c639-4f4c-86e1-b4ba6dab5bec)
(Das im GIF angezeigte Repository sind alle öffentliche Repositories der Organisation, der ich angehöre.)

gh-otui ist ein CLI-Tool, das gh mit ghq und fuzzy finder (peco, fzf) kombiniert.  
Es ermöglicht das Durchsuchen und Anzeigen von Repositories in der Organisation mithilfe von fuzzy finder und das Klonen mit ghq. Besonders praktisch ist es, wenn man an mehreren Repositories arbeitet und den Reponamen kennt, da man nur mit der CLI das Klonen abschließen kann.

## Funktionen

- Anzeige der Listen von GitHub-Organisations-Repositories
- Interaktive Repository-Auswahl mithilfe von fuzzy finder
- Klonen des ausgewählten Repositories mit ghq (bei nicht geklonten Repositories)
- Visuelle Anzeige der geklonten Repositories (✓-Markierung)

## Vorrausgesetzte Tools

- [GitHub CLI](https://cli.github.com/) (gh)
- [ghq](https://github.com/x-motemen/ghq)
- [peco](https://github.com/peco/peco)
  - oder [fzf](https://github.com/junegunn/fzf). Durch Setzen der Umgebungsvariable `GH_OTUI_SELECTOR` auf `fzf` kann fzf verwendet werden. Falls keine Umgebungsvariable angegeben ist, wird das installierte Tool (entweder peco oder fzf) verwendet. Wenn beide installiert sind, hat peco Vorrang.
  
## Installation

```bash
gh extension install n3xem/gh-otui
```

## Benutzung

1. Erstellen Sie den Cache der Repositories Ihrer Organisation:

```bash
gh otui --cache
```

Der Cache wird in `~/.config/gh/extensions/gh-otui/cache.json` gespeichert.

2. Führen Sie folgenden Befehl aus:

```bash
gh otui
```

3. Wählen Sie das gewünschte Repository im fuzzy finder Interface aus
   - Die ✓-Markierung zeigt an, dass das Repository bereits geklont wurde
   - Wenn ein nicht geklontes Repository ausgewählt wird, erfolgt das Klonen mit ghq
   - Der Prüfvorgang, ob ein Repository geklont wurde, erfolgt durch die Überprüfung des Pfades von `ghq root`

4. Der lokale Pfad des ausgewählten Repositories wird auf der Standardausgabe ausgegeben.
   - In Verbindung mit dem cd-Befehl ist dies sehr praktisch für einen schnellen Wechsel.
   - Beispiel: `cd $(gh otui)`

## Ausgabeformat

Die Repositories werden im folgenden Format angezeigt:

- ✓: Markierung für geklonte Repositories
- organization-name: Name der GitHub-Organisation
- repository-name: Name des Repositories