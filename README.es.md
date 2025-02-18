# gh-otui

/oˈtuː.i/ se lee así.

gh-otui = gh + org + tui

![otui](https://github.com/user-attachments/assets/0c7626eb-c639-4f4c-86e1-b4ba6dab5bec)  
(Todos los repositorios mostrados en GIF son públicos de la organización a la que pertenezco)

gh-otui es una herramienta CLI que combina gh, ghq y un buscador difuso (peco, fzf).  
Permite buscar y explorar repositorios de la organización utilizando la mecánica de un buscador difuso, así como clonarlos con ghq. Es especialmente útil cuando se está desarrollando a través de múltiples repositorios, ya que si se conoce el nombre del repositorio, se puede completar el clonaje solo con la CLI.

## Funcionalidades

- Visualización de la lista de repositorios de la organización en GitHub
- Selección interactiva de repositorios utilizando un buscador difuso
- Clonación del repositorio seleccionado mediante ghq (en caso de que no esté clonado)
- Visualización visual de los repositorios clonados (✓ marca)

## Herramientas requeridas

- [GitHub CLI](https://cli.github.com/) (gh)
- [ghq](https://github.com/x-motemen/ghq)
- [peco](https://github.com/peco/peco)
  - O [fzf](https://github.com/junegunn/fzf). Se puede utilizar fzf configurando la variable de entorno `GH_OTUI_SELECTOR` a `fzf`. Si no se especifica ninguna variable de entorno, se usará el que tenga instalada, ya sea peco o fzf. Si ambos están instalados, se priorizará peco.
  
## Instalación

```bash
gh extension install n3xem/gh-otui
```

## Uso

1. Crea un caché de los repositorios de la organización a la que perteneces:

```bash
gh otui --cache
```

El caché se guardará en `~/.config/gh/extensions/gh-otui/cache.json`.

2. Ejecuta el siguiente comando:

```bash
gh otui
```

3. Selecciona el repositorio deseado en la interfaz del buscador difuso
   - La marca ✓ indica los repositorios que ya han sido clonados
   - Al seleccionar un repositorio no clonado, se procederá a clonarlo mediante ghq
   - La determinación de clonado se realiza verificando la ruta de `ghq root`

4. La ruta local del repositorio seleccionado se mostrará en la salida estándar.
   - Es útil usarlo en combinación con el comando cd para poder moverse rápidamente.
   - Ejemplo: `cd $(gh otui)`

## Formato de salida

Los repositorios se mostrarán en el siguiente formato:

- ✓: Marca que indica repositorios clonados
- organization-name: Nombre de la organización en GitHub
- repository-name: Nombre del repositorio