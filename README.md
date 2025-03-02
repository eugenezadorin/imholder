# ImHolder

**ImHolder** is a self-hosted image placeholder generator. It dynamically creates images with customizable size, background color, text, and format. The service is designed to be simple, flexible, and easy to integrate into your development workflow.

> **Disclaimer**: This project was developed with the assistance of AI in a "vibe coding" style. Use it at your own risk. The developers are not responsible for any issues that may arise from its use.

---

## Features

- **Dynamic image generation**: Create images with custom dimensions, background colors, and text.
- **Multiple formats**: Supports PNG, JPG, and SVG formats.
- **Customizable colors**: Use predefined colors or specify custom hex codes.
- **Text customization**: Add custom text or use the default text (image dimensions).
- **Network delay simulation**: Simulate network delays for testing purposes.

## Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/eugenezadorin/imholder.git
   cd imholder
   ```

2. **Build the application**:
   ```bash
   go build -o imholder .
   ```

3. **Run the server**:
   ```bash
   ./imholder
   ```
   By default, the server runs on port `8004`. You can change the port using the `IMHOLDER_PORT` environment variable or a command-line flag `-p`.

## Usage

### Endpoint

The service is accessible via the following URL pattern:

```
/{width}x{height}.{format}?bg={background}&text={text}&text_color={text_color}&delay={delay}
```

- `{width}`: Width of the image (e.g., `300`).
- `{height}`: Height of the image (e.g., `200`).
- `{format}`: Image format (`png`, `jpg`, or `svg`). Defaults to `png`.
- `bg`: Background color (predefined color name or hex code). Defaults to light gray.
- `text`: Custom text to display on the image. Defaults to the image dimensions (e.g., `300x200`).
- `text_color`: Text color (predefined color name or hex code). Defaults to dark gray.
- `delay`: Simulate network delay in milliseconds. Can be a single value (e.g., `500`) or a random value from range (e.g., `200-1000`).

### Examples

1. **Basic image**:
   ```
   /300x200.png
   ```
   Generates a 300x200 PNG image with a light gray background and the text "300x200" in dark gray.

2. **Custom background and text**:
   ```
   /400x300.jpg?bg=blue&text=Hello%20World&text_color=white
   ```
   Generates a 400x300 JPG image with a blue background and the text "Hello World" in white.

3. **Custom hex colors**:
   ```
   /500x500.png?bg=ffcc00&text=Custom%20Colors&text_color=003366
   ```
   Generates a 500x500 PNG image with a yellow background (`#ffcc00`) and the text "Custom Colors" in dark blue (`#003366`).

4. **Simulate network delay**:
   ```
   /600x400.svg?delay=1000
   ```
   Generates a 600x400 SVG image after a 1-second delay.

5. **Random delay within a range**:
   ```
   /800x600.jpg?delay=200-1000
   ```
   Generates an 800x600 JPG image after a random delay between 200ms and 1000ms.

## Predefined Colors

The following predefined colors are available for the `bg` and `text_color` parameters:

| Color Name | Hex Code  |
|------------|-----------|
| red        | `#E63946` |
| orange     | `#FFA500` |
| yellow     | `#FFC857` |
| green      | `#3CB371` |
| blue       | `#1E90FF` |
| purple     | `#9370DB` |
| pink       | `#FFB6C1` |
| brown      | `#8B4513` |
| gray       | `#808080` |
| lightgray  | `#D3D3D3` |
| darkgray   | `#404040` |

## Environment Variables

The following environment variables can be used to configure the application:

| Variable         | Description                          | Default Value |
|------------------|--------------------------------------|---------------|
| `IMHOLDER_PORT`  | Port on which the server will run.   | `8004`        |

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.