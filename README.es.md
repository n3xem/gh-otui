# gh-otui

/oˈtuː.i/ se lee así.

gh-otui = gh + org + tui

![otui](https://github.com/user-attachments/assets/0c7626eb-c639-4f4c-86e1-b4ba6dab5bec)
(Los repositorios mostrados en el GIF son todos públicos de la organización a la que pertenezco)

gh-otui es una herramienta de línea de comandos (CLI) que combina gh, ghq y un buscador difuso (peco, fzf).  
Permite buscar y navegar a través de las organizaciones y repositorios propios utilizando el sistema de buscador difuso, y clonar con ghq. Especialmente útil cuando desarrollas a través de múltiples repositorios, ya que si conoces el nombre del repositorio, puedes completar el clonado solo con la CLI.

## Funcionalidades

- Visualización de la lista de organizaciones de GitHub y de los propios repositorios
- Selección interactiva de repositorios usando un buscador difuso
- Clonación de repositorios seleccionados con ghq (si no han sido clonado previamente)
- Visualización visual de los repositorios clonados (✓ marca)

## Herramientas requeridas

- [GitHub CLI](https://cli.github.com/) (gh)
- [ghq](https://github.com/x-motemen/ghq)
- [peco](https://github.com/peco/peco)
  - O [fzf](https://github.com/junegunn/fzf). Puedes usar fzf configurando la variable de entorno `GH_OTUI_SELECTOR` a `fzf`. Si no hay especificación de variable de entorno, se utilizará la que tenga instalada entre peco y fzf. Si ambas están instaladas, se dará prioridad a peco.
  
## Instalación

```bash
gh extension install n3xem/gh-otui
```

## Uso

1. Simplemente ejecuta el comando gh otui. En la primera ejecución se creará una caché que almacena la lista de repositorios a obtener.

```bash
gh otui
```

2. Selecciona el repositorio deseado en la interfaz del buscador difuso.
   - ✓ marca indica los repositorios que ya han sido clonados.
   - Al seleccionar un repositorio no clonado, se realizará la clonación mediante ghq.
   - La determinación de si un repositorio está clonado se verifica comprobando la ruta de `ghq root`.

3. La ruta local del repositorio seleccionado se mostrará en la salida estándar.
   - Es conveniente usarlo en combinación con el comando cd para poder moverte rápidamente.
   - Ejemplo: `cd $(gh otui)`

## Formato de salida

Los repositorios se mostrarán en el siguiente formato:

- ✓: Marca que indica un repositorio clonado
- organization-name: Nombre de la organización en GitHub
- repository-name: Nombre del repositorio

## Acerca de la caché

gh-otui utiliza la siguiente estructura de caché:

- **Ubicación de la caché**: `~/.config/gh/extensions/gh-otui/`
- **Duración**: 1 hora (se considera que está obsoleto tras 1 hora desde la última actualización, y se actualiza automáticamente en segundo plano)
- **Archivo de metadatos**: `_md.json` - Guarda la última hora de actualización de la caché
- **Directorio de host**: Se crea un directorio para cada host de GitHub (ej: `github.com`)
- **Archivo de organización**: Se crea un archivo `{organization}.json` para cada organización. Guarda información sobre los repositorios.

La caché se actualizará en los siguientes casos:
1. La primera ejecución (si la caché no existe)
2. Si se ha vencido el tiempo de validez de la caché (1 hora) (actualización automática en segundo plano)

Eliminación de la caché: Puedes eliminar el directorio de caché con el comando `gh otui clear`.