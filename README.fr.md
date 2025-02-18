# gh-otui

/oˈtuː.i/ se lit.

gh-otui = gh + org + tui

![otui](https://github.com/user-attachments/assets/0c7626eb-c639-4f4c-86e1-b4ba6dab5bec)
(Tous les dépôts affichés en GIF sont publics de l'organisation à laquelle j'appartiens)

gh-otui est un outil CLI combinant gh, ghq et un fuzzy finder (peco, fzf).  
Il permet d'explorer et de visualiser les dépôts d'une organisation en utilisant le mécanisme du fuzzy finder, et de les cloner avec ghq. En particulier, lorsqu'on développe à travers plusieurs dépôts, il est pratique de pouvoir cloner directement via la CLI si l'on connaît simplement le nom du dépôt.

## Fonctionnalités

- Affichage de la liste des dépôts d'organisation sur GitHub
- Sélection interactive des dépôts à l'aide d'un fuzzy finder
- Clonage via ghq du dépôt sélectionné (si non cloné)
- Affichage visuel des dépôts clonés (✓)

## Outils Pré-requis

- [GitHub CLI](https://cli.github.com/) (gh)
- [ghq](https://github.com/x-motemen/ghq)
- [peco](https://github.com/peco/peco)
  - ou [fzf](https://github.com/junegunn/fzf). En configurant la variable d'environnement `GH_OTUI_SELECTOR` à `fzf`, vous pouvez utiliser fzf. Si aucune variable d'environnement n'est spécifiée, l'outil installé (peco ou fzf) sera utilisé. S'ils sont tous deux installés, peco sera priorisé.
  
## Installation

```bash
gh extension install n3xem/gh-otui
```

## Utilisation

1. Créez le cache des dépôts de l'organisation à laquelle vous appartenez :

```bash
gh otui --cache
```

Le cache sera sauvegardé dans `~/.config/gh/extensions/gh-otui/cache.json`.

2. Exécutez la commande suivante :

```bash
gh otui
```

3. Sélectionnez le dépôt souhaité dans l'interface du fuzzy finder
   - Le marqueur ✓ indique les dépôts déjà clonés
   - Si vous sélectionnez un dépôt non cloné, cela déclenchera un clonage via ghq
   - La détermination des dépôts clonés est effectuée en vérifiant le chemin de `ghq root`

4. Le chemin local du dépôt sélectionné sera affiché en sortie standard.
   - Il est pratique de l'utiliser en conjonction avec la commande cd pour y accéder immédiatement.
   - Exemple : `cd $(gh otui)`

## Format de sortie

Les dépôts sont affichés dans le format suivant :

- ✓ : Marqueur indiquant un dépôt cloné
- organization-name : Nom de l'organisation GitHub
- repository-name : Nom du dépôt