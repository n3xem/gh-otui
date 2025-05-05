# gh-otui

/oˈtuː.i/ wird gelesen.

gh-otui = gh + org + tui

![otui](https://github.com/user-attachments/assets/0c7626eb-c639-4f4c-86e1-b4ba6dab5bec)  
(Das im GIF angezeigte Repository sind alles öffentliche Repositories der Organisation, der ich angehöre.)

gh-otui ist ein CLI-Tool, das gh, ghq und einen fuzzy finder (peco, fzf) kombiniert.  
Es ermöglicht das Durchsuchen und Anzeigen von Organisationen und eigenen Repositories mithilfe der fuzzy finder-Methode und das Klonen über ghq. Besonders nützlich ist es, wenn man an mehreren Repositories arbeitet, da man nur den Repositoriennamen kennen muss, um über die CLI den Klonvorgang abzuschließen.

## Funktionen

- Anzeige der Liste von GitHub-Organisationen und eigenen Repositories
- Interaktive Repository-Auswahl mit einem fuzzy finder
- Klonen des ausgewählten Repositories mit ghq (wenn noch nicht geklont)
- Visuelle Anzeige von bereits geklonten Repositories (✓-Markierung)

## Voraussetzungen

- [GitHub CLI](https://cli.github.com/) (gh)
- [ghq](https://github.com/x-motemen/ghq)
- [peco](https://github.com/peco/peco)
  - Oder [fzf](https://github.com/junegunn/fzf). Sie können fzf verwenden, indem Sie die Umgebungsvariable `GH_OTUI_SELECTOR` auf `fzf` setzen. Wenn keine Umgebungsvariable angegeben ist, wird die installierte Anwendung zwischen peco und fzf verwendet. Wenn beide installiert sind, hat peco Vorrang.

## Installation

```bash
gh extension install n3xem/gh-otui
```

## Verwendung

1. Führen Sie einfach den Befehl `gh otui` aus. Beim ersten Mal wird ein Cache erstellt, der die Liste der zu erhaltenden Repositories speichert.

```bash
gh otui
```

2. Wählen Sie im fuzzy finder-Interface das gewünschte Repository aus.
   - Das ✓-Symbol zeigt an, dass das Repository bereits geklont wurde.
   - Wenn ein noch nicht geklontes Repository ausgewählt wird, erfolgt das Klonen mit ghq.
   - Die Feststellung, ob ein Repository geklont wurde, erfolgt durch Überprüfung des Pfads von `ghq root`.

3. Der lokale Pfad des ausgewählten Repositories wird in der Standardausgabe angezeigt.
   - Es ist praktisch, wenn Sie es zusammen mit dem cd-Befehl verwenden, da Sie sofort wechseln können.
   - Beispiel: `cd $(gh otui)`

## Ausgabeformat

Die Repositories werden im folgenden Format angezeigt:

- ✓: Markierung für geklonte Repositories
- organization-name: Name der GitHub-Organisation
- repository-name: Name des Repositories

## Cache-Informationen

gh-otui verwendet eine Cache-Struktur wie folgt:

- **Cache-Speicherort**: `~/.config/gh/extensions/gh-otui/`
- **Gültigkeitsdauer**: 1 Stunde (Nach 1 Stunde ab letzter Aktualisierung wird dies als alt betrachtet und im Hintergrund automatisch aktualisiert.)
- **Metadatendatei**: `_md.json` - Speichert den letzten Aktualisierungszeitpunkt des Caches
- **Hostverzeichnis**: Für jeden GitHub-Host (z. B.: `github.com`) wird ein Verzeichnis erstellt.
- **Organisationsdatei**: Für jede Organisation wird eine `{organisation}.json` Datei erstellt, die die Repository-Informationen speichert.

Der Cache wird in folgenden Fällen aktualisiert:
1. Bei der ersten Ausführung (wenn der Cache nicht vorhanden ist)
2. Wenn die Cache-Gültigkeitsdauer (1 Stunde) abgelaufen ist (wird im Hintergrund automatisch aktualisiert)

Cache löschen: Mit dem Befehl `gh otui clear` können Sie das Cache-Verzeichnis löschen.