# gh-otui

/oˈtuː.i/ se lit.

gh-otui = gh + org + tui

![otui](https://github.com/user-attachments/assets/0c7626eb-c639-4f4c-86e1-b4ba6dab5bec)
(Les dépôts affichés en GIF sont tous publics de l'organisation à laquelle j'appartiens.)

gh-otui est un outil CLI qui combine gh, ghq et un fuzzy finder (peco, fzf).  
Il permet de rechercher et de naviguer à travers les organisations et ses propres dépôts en utilisant le mécanisme de fuzzy finder, et de cloner via ghq. C'est particulièrement pratique si vous développez sur plusieurs dépôts, car tant que vous connaissez le nom du dépôt, vous pouvez cloner entièrement via la CLI.

## Fonctionnalités

- Affichage de la liste des organisations GitHub et de vos dépôts
- Sélection interactive de dépôts utilisant un fuzzy finder
- Clonage d'un dépôt sélectionné avec ghq (si non cloné)
- Affichage visuel des dépôts clonés (✓)

## Outils requis

- [GitHub CLI](https://cli.github.com/) (gh)
- [ghq](https://github.com/x-motemen/ghq)
- [peco](https://github.com/peco/peco)
  - Ou [fzf](https://github.com/junegunn/fzf). En définissant la variable d'environnement `GH_OTUI_SELECTOR` à `fzf`, vous pouvez utiliser fzf. En l'absence de cette spécification de variable, l'outil installé entre peco et fzf sera utilisé. Si les deux sont installés, peco a la priorité.

## Installation

```bash
gh extension install n3xem/gh-otui
```

## Utilisation

1. Exécutez simplement la commande `gh otui`. Lors de la première exécution, un cache contenant la liste des dépôts à obtenir sera créé.

```bash
gh otui
```

2. Sélectionnez le dépôt souhaité dans l'interface fuzzy finder
   - Le marqueur ✓ indique les dépôts déjà clonés
   - Si vous sélectionnez un dépôt non cloné, il sera cloné via ghq
   - La détermination des dépôts clonés se fait en vérifiant le chemin de `ghq root`

3. Le chemin local du dépôt sélectionné sera affiché en sortie standard.
   - Cela est pratique à utiliser avec la commande `cd`, car vous pouvez y accéder directement.
   - Exemple : `cd $(gh otui)`

## Format de sortie

Les dépôts sont affichés au format suivant :

- ✓ : Marque indiquant un dépôt cloné
- organization-name : Nom de l'organisation GitHub
- repository-name : Nom du dépôt

## À propos du cache

gh-otui utilise une structure de cache comme suit :

- **Emplacement du cache** : `~/.config/gh/extensions/gh-otui/`
- **Durée de validité** : 1 heure (si plus d'une heure s'est écoulée depuis la dernière mise à jour, il sera considéré comme obsolète et sera automatiquement mis à jour en arrière-plan)
- **Fichier de métadonnées** : `_md.json` - Enregistre l'heure de la dernière mise à jour du cache
- **Répertoire hôte** : Création d'un répertoire pour chaque hôte GitHub (ex : `github.com`)
- **Fichiers d'organisation** : Création d'un fichier `{organization}.json` pour chaque organisation. Enregistre les informations sur les dépôts

Le cache est mis à jour dans les cas suivants :
1. Lors de l'exécution initiale (si le cache n'existe pas)
2. Lorsque la durée de validité du cache (1 heure) est expirée (mise à jour automatique en arrière-plan)

Suppression du cache : Vous pouvez supprimer le répertoire du cache avec la commande `gh otui clear`.