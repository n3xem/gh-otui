# gh-otui

/seɪˈtuː.i/ se lee así.

gh-otui = gh + org + tui

![otui](https://github.com/user-attachments/assets/0c7626eb-c639-4f4c-86e1-b4ba6dab5bec)
(Los repositorios mostrados en GIF son todos de organizaciones públicas a las que pertenezco)

gh-otui es una herramienta de línea de comandos (CLI) que combina gh, ghq y buscadores difusos (peco, fzf).  
Permite buscar y navegar a través de los repositorios de la organización utilizando la funcionalidad de buscador difuso, y clonarlos con ghq. Especialmente útil cuando se trabaja en múltiples repositorios, ya que si conoces el nombre del repositorio, puedes completar el proceso de clonación solo con la CLI.

## Funciones

- Visualización de la lista de repositorios de la organización de GitHub
- Selección interactiva de repositorios usando un buscador difuso
- Clonación del repositorio seleccionado mediante ghq (en caso de no estar clonado)
- Visualización visual de los repositorios clonados (marca ✓)

## Herramientas requeridas

- [GitHub CLI](https://cli.github.com/) (gh)
- [ghq](https://github.com/x-motemen/ghq)
- [peco](https://github.com/peco/peco)
  - O [fzf](https://github.com/junegunn/fzf). Puedes usar fzf configurando la variable de entorno `GH_OTUI_SELECTOR` a `fzf`. Si no se especifica una variable de entorno, se utilizará la herramienta que esté instalada, ya sea peco o fzf. Si ambas están instaladas, se dará prioridad a peco.
  
## Instalación

```bash
gh extension install n3xem/gh-otui
```

## Uso

1. Crea la caché de los repositorios de la organización a la que perteneces:

```bash
gh otui --cache
```

La caché se guardará en `~/.config/gh/extensions/gh-otui/cache.json`.

2. Ejecuta el siguiente comando:

```bash
gh otui
```

3. Selecciona el repositorio deseado en la interfaz del buscador difuso.
   - La marca ✓ indica que el repositorio ya ha sido clonado.
   - Si seleccionas un repositorio no clonado, se llevará a cabo la clonación a través de ghq.
   - La determinación de si está clonado se realiza verificando la ruta de `ghq root`.

4. La ruta local del repositorio seleccionado se mostrará en la salida estándar.
   - Es conveniente usarlo junto con el comando cd para poder desplazarte rápidamente.
   - Ejemplo: `cd $(gh otui)`

## Formato de salida

Los repositorios se muestran en el siguiente formato:

- ✓: Marca que indica un repositorio clonado
- organization-name: Nombre de la organización de GitHub
- repository-name: Nombre del repositorio