function generateHexagonGrid() {
  const grid: { q: number; r: number; s: number }[] = [];

  // Define rows in terms of q-coordinates for each row to ensure perfect symmetry
  const rows: [number, number[]][] = [
    // r-value, array of q-values for this row
    [-2, [0, 1, 2]], // Top row (3 hexagons)
    [-1, [-1, 0, 1, 2]], // Second row (4 hexagons)
    [0, [-2, -1, 0, 1, 2]], // Middle row (5 hexagons)
    [1, [-2, -1, 0, 1]], // Fourth row (4 hexagons)
    [2, [-2, -1, 0]], // Bottom row (3 hexagons)
  ];

  // Generate the grid based on the defined rows
  for (const [r, qValues] of rows) {
    for (const q of qValues) {
      // Calculate s using the constraint that q + r + s = 0 in cube coordinates
      const s = -q - r;
      grid.push({ q, r, s });
    }
  }

  return grid;
}

/**
 * Converts cube coordinates to pixel positions for a pointy-top hexagon with spacing
 * @param {Object} hex - Hexagon coordinates in cube format {q, r, s}
 * @param {number} size - Size of the hexagon (from center to corner)
 * @param {number} spacing - Spacing between hexagons
 * @returns {Object} Pixel coordinates {x, y}
 */
function hexToPixel(hex, size, spacing) {
  // Calculate the scaling factor to include spacing
  // For a pointy-top hexagon with size = side length, the width is sqrt(3) * size
  // We add the spacing to this width
  const spacingFactor = 1 + spacing / (Math.sqrt(3) * size);

  // Apply the spacing factor to the standard conversion
  const x = spacingFactor * size * (Math.sqrt(3) * hex.q + (Math.sqrt(3) / 2) * hex.r);
  const y = spacingFactor * size * ((3 / 2) * hex.r);

  return { x, y };
}

/**
 * Creates the points string for an SVG hexagon path
 * @param {Object} center - Center point {x, y} of the hexagon
 * @param {number} size - Size of the hexagon (from center to corner)
 * @returns {string} Points string for SVG polygon
 */
function getHexagonPoints(center, size) {
  let points: string[] = [];
  for (let i = 0; i < 6; i++) {
    const angle = ((2 * Math.PI) / 6) * (i + 0.5); // +0.5 to make it pointy-top
    const x = center.x + size * Math.cos(angle);
    const y = center.y + size * Math.sin(angle);
    points.push(`${x},${y}`);
  }
  return points.join(" ");
}

/**
 * Generates the points for a rectangle between hexagons
 * @param {Object} center1 - Center of first hexagon {x, y}
 * @param {Object} center2 - Center of second hexagon {x, y}
 * @param {number} size - Size of each hexagon
 * @param {number} spacing - Spacing between hexagons
 * @returns {string} Points for the rectangle SVG polygon
 */
function getSpacerRectanglePoints(center1, center2, size, spacing) {
  // Calculate direction vector between centers
  const dx = center2.x - center1.x;
  const dy = center2.y - center1.y;
  const distance = Math.sqrt(dx * dx + dy * dy);

  // Normalize direction vector
  const nx = dx / distance;
  const ny = dy / distance;

  // Calculate perpendicular vector (rotate 90 degrees)
  const px = -ny;
  const py = nx;

  // Width of the rectangle is the spacing
  const width = spacing;

  // Get the corner points of the rectangle
  // The rectangle extends from the edge of hex1 to the edge of hex2
  const hex1Edge = {
    x: center1.x + nx * size,
    y: center1.y + ny * size,
  };

  const hex2Edge = {
    x: center2.x - nx * size,
    y: center2.y - ny * size,
  };

  // Calculate the four corners of the rectangle
  const p1 = {
    x: hex1Edge.x + (px * width) / 2,
    y: hex1Edge.y + (py * width) / 2,
  };

  const p2 = {
    x: hex1Edge.x - (px * width) / 2,
    y: hex1Edge.y - (py * width) / 2,
  };

  const p3 = {
    x: hex2Edge.x - (px * width) / 2,
    y: hex2Edge.y - (py * width) / 2,
  };

  const p4 = {
    x: hex2Edge.x + (px * width) / 2,
    y: hex2Edge.y + (py * width) / 2,
  };

  return `${p1.x},${p1.y} ${p2.x},${p2.y} ${p3.x},${p3.y} ${p4.x},${p4.y}`;
}

/**
 * Creates the SVG markup string directly with green spaces between hexagons
 * @param {number} hexSize - Size of each hexagon (side length)
 * @param {number} spacing - Spacing between hexagons in pixels
 * @returns {string} SVG markup as a string
 */
function getHexagonGridSVGMarkup(hexSize = 30, spacing = 6) {
  const grid = generateHexagonGrid();

  // Calculate required SVG size and center offset
  let minX = Infinity,
    maxX = -Infinity,
    minY = Infinity,
    maxY = -Infinity;

  const pixelCoordinates = grid.map((hex) => {
    const pixel = hexToPixel(hex, hexSize, spacing);
    minX = Math.min(minX, pixel.x - hexSize);
    maxX = Math.max(maxX, pixel.x + hexSize);
    minY = Math.min(minY, pixel.y - hexSize);
    maxY = Math.max(maxY, pixel.y + hexSize);
    return { hex, pixel };
  });

  // Calculate SVG dimensions with padding
  const padding = 20;
  const svgWidth = maxX - minX + padding * 2;
  const svgHeight = maxY - minY + padding * 2;

  // Start SVG markup
  let svgMarkup = `<svg width="${svgWidth}" height="${svgHeight}" viewBox="0 0 ${svgWidth} ${svgHeight}" xmlns="http://www.w3.org/2000/svg">`;

  // Group for green spacers
  svgMarkup += `<g id="spacers">`;

  // Find adjacent hexagons to create spacers between them
  const spacers = [];
  for (let i = 0; i < pixelCoordinates.length; i++) {
    for (let j = i + 1; j < pixelCoordinates.length; j++) {
      const hex1 = pixelCoordinates[i].hex;
      const hex2 = pixelCoordinates[j].hex;

      // Check if hexagons are adjacent (distance of 1 in cube coordinates)
      const distance =
        Math.abs(hex1.q - hex2.q) + Math.abs(hex1.r - hex2.r) + Math.abs(hex1.s - hex2.s);
      if (distance === 2) {
        // In cube coordinates, adjacent hexes have a distance of 2
        spacers.push([pixelCoordinates[i], pixelCoordinates[j]]);
      }
    }
  }

  // Add green spacers between adjacent hexagons
  spacers.forEach(([item1, item2]) => {
    const center1 = {
      x: item1.pixel.x - minX + padding,
      y: item1.pixel.y - minY + padding,
    };

    const center2 = {
      x: item2.pixel.x - minX + padding,
      y: item2.pixel.y - minY + padding,
    };

    const points = getSpacerRectanglePoints(center1, center2, hexSize, spacing);

    svgMarkup += `<polygon points="${points}" fill="#4CAF50" />`;
  });

  svgMarkup += `</g>`;

  // Group for hexagons (drawn on top of spacers)
  svgMarkup += `<g id="hexagons">`;

  // Add each hexagon to the SVG
  pixelCoordinates.forEach((item) => {
    const { hex, pixel } = item;
    // Center the grid in the SVG by applying the offset
    const centeredPixel = {
      x: pixel.x - minX + padding,
      y: pixel.y - minY + padding,
    };

    // Generate points for the hexagon
    const points = getHexagonPoints(centeredPixel, hexSize);

    // Determine styles based on coordinates
    let fillColor = "#eaeaea";
    if (hex.q === 0 && hex.r === 0 && hex.s === 0) {
      fillColor = "#ff9999"; // Color the center hexagon differently
    }

    // Add polygon element to SVG markup
    svgMarkup += `<polygon points="${points}" fill="${fillColor}" stroke="#333" stroke-width="1">`;
    svgMarkup += `<title>q: ${hex.q}, r: ${hex.r}, s: ${hex.s}</title>`;
    svgMarkup += `</polygon>`;
  });

  svgMarkup += `</g>`;

  // Close SVG markup
  svgMarkup += `</svg>`;

  return svgMarkup;
}

/**
 * Renders the hexagon grid SVG to the DOM
 * @param {string} containerId - ID of the container element
 * @param {number} hexSize - Size of each hexagon
 * @param {number} spacing - Spacing between hexagons in pixels
 */
function renderHexagonGridSVG(containerId, hexSize = 30, spacing = 6) {
  const container = document.getElementById(containerId);
  if (!container) return;

  container.innerHTML = getHexagonGridSVGMarkup(hexSize, spacing);
}

// For direct use:
window.onload = function () {
  // Create container div if it doesn't exist
  let container = document.getElementById("hexagonGridContainer");
  if (!container) {
    container = document.createElement("div");
    container.id = "hexagonGridContainer";
    document.body.appendChild(container);
  }

  // Render the hexagon grid SVG with hexSize=30 and spacing=6
  renderHexagonGridSVG("hexagonGridContainer", 30, 6);
};

const ports = {
  1: {
    multiX: 2.6,
    multiY: -1.25,
    multiH: 1,
    multiW: 0.5,
    index: 3,
  },
  4: {
    multiX: -0.775,
    multiY: 2.125,
    multiH: 1,
    multiW: 0.5,
    index: 3,
  },
  10: {
    multiX: 2.6,
    multiY: 1,
    multiH: 1,
    multiW: 0.5,
    index: 5,
  },
  13: {
    multiX: 2.6,
    multiY: 1,
    multiH: 1,
    multiW: 0.5,
    index: 2,
  },
  14: {
    multiX: 1,
    multiY: -1.6,
    multiH: 0.5,
    multiW: 1,
    index: 4,
  },
  18: {
    multiX: 2.6,
    multiY: 1,
    multiH: 1,
    multiW: 0.5,
    index: 5,
  },
  24: {
    multiX: 1,
    multiY: -1.6,
    multiH: 0.5,
    multiW: 1,
    index: 1,
  },
  25: {
    multiX: -1.6,
    multiY: 1,
    multiH: 1,
    multiW: 0.5,
    index: 3,
  },
  26: {
    multiX: 2.6,
    multiY: 1,
    multiH: 1,
    multiW: 0.5,
    index: 0,
  },
  29: {
    multiX: 1,
    multiY: 2.6,
    multiH: 0.5,
    multiW: 1,
    index: 4,
  },
  36: {
    multiX: -1.6,
    multiY: 1,
    multiH: 1,
    multiW: 0.5,
    index: 0,
  },
  39: {
    multiX: -1.6,
    multiY: 1,
    multiH: 1,
    multiW: 0.5,
    index: 5,
  },
  40: {
    multiX: 1,
    multiY: 2.6,
    multiH: 0.5,
    multiW: 1,
    index: 1,
  },
  47: {
    multiX: -0.775,
    multiY: -1.6,
    multiH: 0.5,
    multiW: 1,
    index: 4,
  },
  50: {
    multiX: -1.6,
    multiY: -0.775,
    multiH: 1,
    multiW: 0.5,
    index: 5,
  },
  51: {
    multiX: 1,
    multiY: 2.6,
    multiH: 0.5,
    multiW: 1,
    index: 1,
  },
  53: {
    multiX: -1.6,
    multiY: 1,
    multiH: 1,
    multiW: 0.5,
    index: 0,
  },
  54: {
    multiX: -1.6,
    multiY: 1,
    multiH: 1,
    multiW: 0.5,
    index: 2,
  },
};
