# gh-otui

/seɪˈtuː.i/ se lee así.

gh-otui = gh + org + tui

![otui](https://github.com/user-attachments/assets/0c7626eb-c639-4f4c-86e1-b4ba6dab5bec)  
(Todos los repositorios mostrados en el GIF son públicos y pertenecen a la organización a la que estoy afiliado.)

gh-otui es una herramienta CLI que combina gh, ghq y un buscador difuso (peco, fzf).  
Permite buscar y explorar organizaciones o los repositorios de uno mismo utilizando un buscador difuso, y clonar con ghq. Es especialmente útil cuando se desarrollan múltiples repositorios, ya que si se conoce el nombre del repositorio, se puede completar el clonaje solo con la CLI.

## Funciones

- Visualización de la lista de organizaciones y repositorios de GitHub
- Selección interactiva de repositorios utilizando un buscador difuso
- Clonación del repositorio seleccionado con ghq (si no ha sido clonado)
- Visualización visual de los repositorios clonados (✓ marca)

## Herramientas necesarias

- [GitHub CLI](https://cli.github.com/) (gh)
- [ghq](https://github.com/x-motemen/ghq)
- [peco](https://github.com/peco/peco) 
  - o [fzf](https://github.com/junegunn/fzf). Puedes utilizar fzf estableciendo la variable de entorno `GH_OTUI_SELECTOR` a `fzf`. Si no se especifica ninguna variable de entorno, se utilizará el que tenga instalado entre peco y fzf. Si ambos están instalados, se priorizará peco.

## Instalación

```bash
gh extension install n3xem/gh-otui
```

## Cómo usar

1. Crea un caché de los repositorios de la organización a la que perteneces:

```bash
gh otui --cache
```

El caché se guardará en `~/.config/gh/extensions/gh-otui/cache.json`.

2. Ejecuta el siguiente comando:

```bash
gh otui
```

3. Selecciona el repositorio deseado en la interfaz del buscador difuso:
   - ✓ La marca indica que el repositorio ya ha sido clonado.
   - Al seleccionar un repositorio no clonado, se realizará el clonaje con ghq.
   - La determinación de si ya ha sido clonado se verifica revisando la ruta de `ghq root`.

4. La ruta local del repositorio seleccionado se mostrará en la salida estándar.
   - Es conveniente usarlo en combinación con el comando cd para moverse rápidamente.
   - Ejemplo: `cd $(gh otui)`

## Formato de salida

Los repositorios se mostrarán en el siguiente formato:

- ✓: Marca que indica un repositorio clonado
- organization-name: Nombre de la organización en GitHub
- repository-name: Nombre del repositorio