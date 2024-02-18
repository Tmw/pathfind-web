(function main(window, document) {
  let canvas, ctx;

  function resizeCanvas(state) {
    canvas.width = state.tile_size * state.grid.width + 1;
    canvas.height = state.tile_size * state.grid.height + 1;
  }

  function mouseMove({ offsetX, offsetY }) {
    window.bridge.mouseMove(
      JSON.stringify({
        x: offsetX,
        y: offsetY,
      })
    );
  }

  function mouseDown() {
    window.bridge.mouseDown();
  }

  function mouseUp() {
    window.bridge.mouseUp();
  }

  function drawStart(x, y, tileSize) {
    ctx.fillStyle = "#090";
    ctx.fillRect(x * tileSize, y * tileSize, tileSize, tileSize);
  }

  function drawFinish(x, y, tileSize) {
    ctx.fillStyle = "#800";
    ctx.fillRect(x * tileSize, y * tileSize, tileSize, tileSize);
  }

  function drawWall(x, y, tileSize) {
    ctx.fillStyle = "#777";
    ctx.fillRect(x * tileSize, y * tileSize, tileSize, tileSize);
  }

  function drawPath(x, y, tileSize) {
    ctx.fillStyle = "#ffcc00";
    ctx.fillRect(x * tileSize, y * tileSize, tileSize, tileSize);
  }

  function drawSelector(x, y, tileSize) {
    ctx.beginPath();
    ctx.strokeStyle = "teal";
    ctx.rect(x * tileSize, y * tileSize, tileSize, tileSize);
    ctx.closePath();
    ctx.stroke();
  }

  function drawOpen(x, y, tileSize) {
    ctx.strokeStyle = "#bbb";
    ctx.beginPath();
    ctx.rect(x * tileSize, y * tileSize, tileSize, tileSize);
    ctx.stroke();
    ctx.closePath();
  }

  function drawGrid(state) {
    ctx.lineWidth = 1;

    // use logic from state.
    const renderFns = {
      open: drawOpen,
      wall: drawWall,
      start: drawStart,
      finish: drawFinish,
      path: drawPath,
    };

    for (let y = 0; y < state.grid.height; y++) {
      for (let x = 0; x < state.grid.width; x++) {
        const fn = renderFns[state.grid.tiles[y][x]];
        if (fn) {
          fn(x, y, state.tile_size);
        }
      }
    }
  }

  function drawCursor(state) {
    const { cursor_position, cursor_tile, tile_size } = state;

    switch (cursor_tile) {
      case "selector":
        drawSelector(cursor_position.x, cursor_position.y, tile_size);
        break;

      case "start":
        drawStart(cursor_position.x, cursor_position.y, tile_size);
        break;

      case "finish":
        drawFinish(cursor_position.x, cursor_position.y, tile_size);
        break;
    }

    drawSelector(cursor_position.x, cursor_position.y, tile_size);
  }

  function draw(state) {
    ctx.reset();
    drawGrid(state);
    drawCursor(state);
  }

  function loop() {
    var state;
    try {
      state = JSON.parse(bridge.getState());
    } catch (e) {
      console.error("Error parsing JSON response from WASM", e);
    }

    resizeCanvas(state);
    draw(state);
    window.requestAnimationFrame(loop);
  }

  function solve() {
    window.bridge.solve();
  }

  function reset() {
    window.bridge.reset();
  }

  function clearPath() {
    window.bridge.clearPath();
  }

  function updateGrid(e) {
    e.preventDefault();
    const tileSize = parseInt(
      document.querySelector("input[name=tileSize]").value
    );
    const height = parseInt(
      document.querySelector("input[name=gridHeight]").value
    );
    const width = parseInt(
      document.querySelector("input[name=gridWidth]").value
    );

    const req = {
      tile_size: tileSize,
      width,
      height,
    };

    bridge.resize(JSON.stringify(req));
  }

  function ready() {
    loadWasm().then(() => {
      // grab references to canvas
      canvas = document.getElementById("canvas");
      canvas.addEventListener("mousemove", mouseMove);
      canvas.addEventListener("mousedown", mouseDown);
      canvas.addEventListener("mouseup", mouseUp);
      ctx = canvas.getContext("2d");

      // grab references to form elements
      const optionsForm = document.getElementById("options-form");
      optionsForm.addEventListener("submit", updateGrid);

      document.getElementById("go-button").addEventListener("click", solve);
      document
        .getElementById("clear-button")
        .addEventListener("click", clearPath);

      document.getElementById("reset-button").addEventListener("click", reset);

      loop();
    });
  }

  async function loadWasm() {
    return new Promise((res) => {
      const go = new window.Go();
      WebAssembly.instantiateStreaming(
        fetch("bridge.wasm"),
        go.importObject
      ).then(({ instance }) => {
        go.run(instance);
        res();
      });
    });
  }

  document.addEventListener("DOMContentLoaded", ready);
})(window, document);
