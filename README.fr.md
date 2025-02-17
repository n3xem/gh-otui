# gh-otui

/seˈtɥiː/ se lit.

gh-otui = gh + org + tui

![otui](https://github.com/user-attachments/assets/0c7626eb-c639-4f4c-86e1-b4ba6dab5bec)  
(L'intégralité des dépôts affichés dans le GIF appartient à l'organisation à laquelle je fais partie.)

gh-otui est un outil CLI combinant gh, ghq et un sélecteur flou (peco, fzf).  
Il permet de parcourir et de rechercher des dépôts d'organisation en utilisant un système de sélection flou, et de cloner avec ghq. Il est particulièrement utile lorsque vous développez à travers plusieurs dépôts car, tant que vous connaissez le nom du dépôt, vous pouvez compléter le clonage uniquement via la CLI.

## Fonctionnalités

- Affichage de la liste des dépôts d'organisation GitHub
- Sélection de dépôts interactive utilisant un sélecteur flou
- Clonage d'un dépôt sélectionné avec ghq (si pas encore cloné)
- Affichage visuel des dépôts clonés (✓)

## Outils Prérequis

- [GitHub CLI](https://cli.github.com/) (gh)
- [ghq](https://github.com/x-motemen/ghq)
- [peco](https://github.com/peco/peco)
  - Ou [fzf](https://github.com/junegunn/fzf). Vous pouvez utiliser fzf en définissant la variable d'environnement `GH_OTUI_SELECTOR` sur `fzf`. En l'absence de spécification de variable d'environnement, il utilisera celui installé entre peco et fzf. Si les deux sont installés, peco aura la priorité.

## Installation

```bash
gh extension install n3xem/gh-otui
```

## Utilisation

1. Créez le cache des dépôts de l'organisation à laquelle vous appartenez :

```bash
gh otui --cache
```

Le cache sera enregistré dans `~/.config/gh/extensions/gh-otui/cache.json`.

2. Exécutez la commande suivante :

```bash
gh otui
```

3. Sélectionnez le dépôt souhaité dans l'interface du sélecteur flou :
   - Le symbole ✓ indique un dépôt déjà cloné.
   - Si vous sélectionnez un dépôt non cloné, un clonage sera effectué avec ghq.
   - La détermination du clonage se fait en vérifiant le chemin de `ghq root`.

4. Le chemin local du dépôt sélectionné sera affiché en sortie standard.
   - Cela est utile lorsque vous l'utilisez avec la commande cd pour vous déplacer rapidement.
   - Exemple : `cd $(gh otui)`

## Format de Sortie

Les dépôts sont affichés dans le format suivant :

- ✓ : Marque pour indiquer un dépôt cloné
- organization-name : Nom de l'organisation GitHub
- repository-name : Nom du dépôt