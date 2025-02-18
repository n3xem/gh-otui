# gh-otui

/oˈtuː.i/ wird gelesen.

gh-otui = gh + org + tui

![otui](https://github.com/user-attachments/assets/0c7626eb-c639-4f4c-86e1-b4ba6dab5bec)
(Das im GIF angezeigte Repository ist alles öffentliche Repositories meiner zugehörigen Organisation.)

gh-otui ist ein CLI-Tool, das gh, ghq und fuzzy finder (peco, fzf) kombiniert.  
Es ermöglicht das Durchsuchen und Anzeigen von Repositories einer Organisation mithilfe des fuzzy finder-Mechanismus und das Klonen mit ghq. Besonders wenn Sie an mehreren Repositories gleichzeitig arbeiten, ist es praktisch, dass Sie, sofern Sie den Repository-Namen kennen, den Klonvorgang nur über die CLI abschließen können.

## Funktionen

- Anzeige einer Liste von GitHub-Organisations-Repositories
- Interaktive Repository-Auswahl mit fuzzy finder
- Klonen des ausgewählten Repositories mit ghq (wenn nicht geklont)
- Visuelle Anzeige geklonter Repositories (✓-Markierung)

## Voraussetzungen

- [GitHub CLI](https://cli.github.com/) (gh)
- [ghq](https://github.com/x-motemen/ghq)
- [peco](https://github.com/peco/peco)
  - oder [fzf](https://github.com/junegunn/fzf). Sie können fzf verwenden, indem Sie die Umgebungsvariable `GH_OTUI_SELECTOR` auf `fzf` setzen. Wenn keine Umgebungsvariable angegeben ist, wird die installierte Option von peco und fzf verwendet. Wenn beide installiert sind, hat peco Vorrang.
  
## Installation

```bash
gh extension install n3xem/gh-otui
```

## Verwendung

1. Erstellen Sie den Cache für die Repositories Ihrer Organisation:

```bash
gh otui --cache
```

Der Cache wird in `~/.config/gh/extensions/gh-otui/cache.json` gespeichert.

2. Führen Sie den folgenden Befehl aus:

```bash
gh otui
```

3. Wählen Sie im fuzzy finder-Interface das gewünschte Repository aus.
   - Das ✓-Zeichen zeigt ein bereits geklontes Repository an.
   - Wenn Sie ein nicht geklontes Repository auswählen, wird es mit ghq geklont.
   - Die Bestimmung, ob geklont wurde, erfolgt durch Überprüfung des Pfades von `ghq root`.

4. Der lokale Pfad des ausgewählten Repositories wird auf die Standardausgabe ausgegeben.
   - Es ist praktisch, es in Kombination mit dem cd-Befehl zu verwenden, um sofort zu wechseln.
   - Beispiel: `cd $(gh otui)`

## Ausgabeformat

Die Repositories werden im folgenden Format angezeigt:

- ✓: Markierung für geklonte Repositories
- organization-name: Name der GitHub-Organisation
- repository-name: Name des Repositories