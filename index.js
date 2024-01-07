(function main(window, document) {
  let canvas, ctx, mouseX, mouseY, isMouseDown, grid;

  let tileSize = 30;
  let gridWidth = 40;
  let gridHeight = 30;

  const tileType = {
    selector: 0,
    open: 1,
    wall: 2,
    start: 3,
    finish: 4,
    candidate: 5,
    visited: 6,
    path: 7,
  };

  let cursorTile = tileType.wall;

  function setGridCell(x, y, cellType) {
    grid[y][x] = cellType;
  }

  function getCellType(x, y) {
    return grid[y][x];
  }

  function resizeCanvas() {
    canvas.width = tileSize * gridWidth + 1;
    canvas.height = tileSize * gridHeight + 1;
  }

  function mouseMove(e) {
    mouseX = e.offsetX;
    mouseY = e.offsetY;

    // support for dragging to place walls
    const { x, y } = mouseToMapCoord();
    if (
      isMouseDown &&
      cursorTile == tileType.wall &&
      getCellType(x, y) === tileType.open
    ) {
      setGridCell(x, y, tileType.wall);
    }

    if (
      isMouseDown &&
      cursorTile == tileType.open &&
      getCellType(x, y) === tileType.wall
    ) {
      setGridCell(x, y, tileType.open);
    }
  }

  function mouseDown() {
    if (isMouseDown) return;
    isMouseDown = true;
    const { x, y } = mouseToMapCoord();
    switch (getCellType(x, y)) {
      case tileType.start:
        setGridCell(x, y, tileType.open);
        cursorTile = tileType.start;
        break;

      case tileType.finish:
        setGridCell(x, y, tileType.open);
        cursorTile = tileType.finish;
        break;

      case tileType.open:
        cursorTile = tileType.wall;
        break;

      case tileType.wall:
        cursorTile = tileType.open;
        break;
    }
  }

  function mouseUp() {
    isMouseDown = false;
    const { x, y } = mouseToMapCoord();
    switch (cursorTile) {
      case tileType.start:
        setGridCell(x, y, tileType.start);
        cursorTile = tileType.selector;
        break;

      case tileType.finish:
        setGridCell(x, y, tileType.finish);
        cursorTile = tileType.selector;
        break;

      case tileType.open:
        setGridCell(x, y, tileType.open);
        cursorTile = tileType.selector;
        break;

      case tileType.wall:
        setGridCell(x, y, tileType.wall);
        cursorTile = tileType.selector;
        break;
    }
  }

  function mouseToMapCoord() {
    return {
      x: Math.min(Math.floor(mouseX / tileSize), gridWidth - 1),
      y: Math.min(Math.floor(mouseY / tileSize), gridHeight - 1),
    };
  }

  function drawStart(x, y) {
    ctx.fillStyle = "#090";
    ctx.fillRect(x * tileSize, y * tileSize, tileSize, tileSize);
  }

  function drawFinish(x, y) {
    ctx.fillStyle = "#800";
    ctx.fillRect(x * tileSize, y * tileSize, tileSize, tileSize);
  }

  function drawWall(x, y) {
    ctx.fillStyle = "#777";
    ctx.fillRect(x * tileSize, y * tileSize, tileSize, tileSize);
  }

  function drawPath(x, y) {
    ctx.fillStyle = "#ffcc00";
    ctx.fillRect(x * tileSize, y * tileSize, tileSize, tileSize);
  }

  function drawSelector(x, y) {
    ctx.beginPath();
    ctx.strokeStyle = "teal";
    ctx.rect(x * tileSize, y * tileSize, tileSize, tileSize);
    ctx.closePath();
    ctx.stroke();
  }

  function drawOpen(x, y) {
    ctx.strokeStyle = "#bbb";
    ctx.beginPath();
    ctx.rect(x * tileSize, y * tileSize, tileSize, tileSize);
    ctx.stroke();
    ctx.closePath();
  }

  function drawGrid() {
    ctx.lineWidth = 1;

    for (let y = 0; y < gridHeight; y++) {
      for (let x = 0; x < gridWidth; x++) {
        switch (grid[y][x]) {
          case tileType.open:
            drawOpen(x, y);
            break;

          case tileType.wall:
            drawWall(x, y);
            break;

          case tileType.start:
            drawStart(x, y);
            break;
          case tileType.path:
            drawPath(x, y);
            break;

          case tileType.finish:
            drawFinish(x, y);
            break;
        }
      }
    }
  }

  function drawCursor() {
    const { x, y } = mouseToMapCoord();
    switch (cursorTile) {
      case tileType.selector:
        drawSelector(x, y);
        break;

      case tileType.start:
        drawStart(x, y);
        break;

      case tileType.finish:
        drawFinish(x, y);
        break;
    }
    drawSelector(x, y);
  }

  function draw() {
    ctx.reset();
    drawGrid();
    drawCursor();
  }

  function loop() {
    draw();
    window.requestAnimationFrame(loop);
  }

  function makeGrid() {
    grid = [];
    for (let y = 0; y < gridHeight; y++) {
      grid[y] = [];
      for (let x = 0; x < gridWidth; x++) {
        grid[y][x] = tileType.open;
      }
    }

    grid[2][4] = tileType.wall;
    grid[6][6] = tileType.start;
    grid[12][18] = tileType.finish;
  }

  function solve() {
    const pathResp = window.bridge.solve(
      JSON.stringify({
        gridWidth,
        gridHeight,
        grid,
      })
    );

    const path = JSON.parse(pathResp);
    // only color the open cells in the path, leaving start and finish as is.
    for (const { x, y } of path) {
      if (getCellType(x, y) === tileType.open) {
        setGridCell(x, y, tileType.path);
      }
    }
  }

  function updateGrid(e) {
    e.preventDefault();
    tileSize = parseInt(document.querySelector("input[name=tileSize]").value);
    gridHeight = parseInt(
      document.querySelector("input[name=gridHeight]").value
    );
    gridWidth = parseInt(document.querySelector("input[name=gridWidth]").value);
    resizeCanvas();
    makeGrid();
  }

  function ready() {
    loadWasm();
    // grab references to canvas
    canvas = document.getElementById("canvas");
    canvas.addEventListener("mousemove", mouseMove);
    canvas.addEventListener("mousedown", mouseDown);
    canvas.addEventListener("mouseup", mouseUp);
    ctx = canvas.getContext("2d");

    // grab references to form elements
    const optionsForm = document.getElementById("options-form");
    optionsForm.addEventListener("submit", updateGrid);

    // listen go clicks on the go button
    document.getElementById("go-button").addEventListener("click", solve);

    // setup everything else
    resizeCanvas();
    makeGrid();
    loop();
  }

  function loadWasm() {
    const go = new window.Go();
    WebAssembly.instantiateStreaming(
      fetch("bridge.wasm"),
      go.importObject
    ).then(({ instance }) => go.run(instance));
  }

  document.addEventListener("DOMContentLoaded", ready);
})(window, document);
