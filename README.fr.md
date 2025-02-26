# gh-otui

/oˈtuː.i/ se lit.

gh-otui = gh + org + tui

![otui](https://github.com/user-attachments/assets/0c7626eb-c639-4f4c-86e1-b4ba6dab5bec)
(Les dépôts affichés dans le GIF sont tous des dépôts publics de mon organisation.)

gh-otui est un outil CLI combinant gh, ghq et un fuzzy finder (peco, fzf).  
Il vous permet de rechercher et de naviguer à travers les organisations et vos propres dépôts en utilisant le système de fuzzy finder, tout en pouvant les cloner avec ghq. C'est particulièrement pratique lorsque vous développez à travers plusieurs dépôts, car tant que vous connaissez le nom du dépôt, vous pouvez cloner uniquement avec la CLI.

## Fonctionnalités

- Affichage de la liste des organisations GitHub et de vos dépôts
- Sélection interactive des dépôts à l'aide d'un fuzzy finder
- Clonage du dépôt sélectionné avec ghq (si non cloné)
- Affichage visuel des dépôts clonés (✓ marque)

## Outils prérequis

- [GitHub CLI](https://cli.github.com/) (gh)
- [ghq](https://github.com/x-motemen/ghq)
- [peco](https://github.com/peco/peco)
  - ou [fzf](https://github.com/junegunn/fzf). Vous pouvez utiliser fzf en définissant la variable d'environnement `GH_OTUI_SELECTOR` sur `fzf`. En l'absence de spécification de variable d'environnement, l'outil installé entre peco et fzf sera utilisé. Si les deux sont installés, peco sera priorisé.
  
## Installation

```bash
gh extension install n3xem/gh-otui
```

## Utilisation

1. Créez le cache des dépôts de votre organisation :

```bash
gh otui --cache
```

Le cache est enregistré dans `~/.config/gh/extensions/gh-otui/cache.json`.

2. Exécutez la commande suivante :

```bash
gh otui
```

3. Sélectionnez le dépôt souhaité dans l'interface du fuzzy finder
   - La marque ✓ indique les dépôts déjà clonés
   - En sélectionnant un dépôt non cloné, le clonage avec ghq sera effectué
   - La détermination du clonage se fait en vérifiant le chemin de `ghq root`

4. Le chemin local du dépôt sélectionné sera affiché dans la sortie standard.
   - Cela est pratique à utiliser en conjonction avec la commande cd, vous permettant de vous déplacer rapidement.
   - Exemple : `cd $(gh otui)`

## Format de sortie

Les dépôts sont affichés dans le format suivant :

- ✓ : Marque indiquant un dépôt cloné
- organization-name : Nom de l'organisation GitHub
- repository-name : Nom du dépôt