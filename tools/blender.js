/**
 * @typedef {Object} OKLABColor
 * @property {number} a
 * @property {number} b
 * @property {number} L
 */

/**
 * @param {OKLABColor} a
 * @param {OKLABColor} b
 * @param {number} fac
 * @returns {OKLABColor}
 */

function mixColors(a, b, fac) {
  return {
    a: mixNumber(a.a % 360, b.a % 360, fac),
    b: mixNumber(a.b % 360, b.b % 360, fac),
    L: mixNumber(a.L, b.L, fac),
  };
}

/**
 * @param {number} a
 * @param {number} b
 * @param {number} fac
 * @returns {number}
 */
function mixNumber(a, b, fac) {
  return a * (1 - fac) + b * fac;
}

/**
 * @param {OKLABColor} a
 * @param {OKLABColor} b
 * @returns {number}
 */
function minimumFactor(a, b) {
  minA = minimumFactorForComponent(a.a, b.a, -0.4, 0.4);
  minB = minimumFactorForComponent(a.b, b.b, -0.4, 0.4);
  minL = minimumFactorForComponent(a.b, b.b, 0, 1);

  return Math.max(minA, minB, minL);
}

/**
 * @param {number} a
 * @param {number} b
 * @param {number} min
 * @param {number} max
 * @returns {number}
 */
function minimumFactorForComponent(a, b, min, max) {
  Math.max((a - b) / (min + a), (a - b) / (max + a));
}

/**
 * @param {OKLABColor} a
 * @param {OKLABColor} b
 * @returns {{Color: OKLABColor, Factor: number }}
 */
function getColorToMixAndFactor(a, b) {
  let Factor = Math.ceil(minimumFactor(a, b) * 1000) / 1000;

  return {
    Color: {
      a: (a.a - b.a * (1 - Factor)) / Factor,
      b: (a.b - b.b * (1 - Factor)) / Factor,
      L: (a.L - b.L * (1 - Factor)) / Factor,
    },
    Factor,
  };
}

// The following lines are copied from earthbound19/OklabExperiments.js
// https://gist.github.com/earthbound19/e7fe15fdf8ca3ef814750a61bc75b5ce#file-oklabexperiments-js

// This gist contains JavaScript functions and tests for:
// - conversion from gamma-corrected (or gamma-compressed) sRGB to linear RGB, to Oklab
// - interpolation through Oklab
// - conversion back to linear RGB, then sRGB
// To use these tests, install nodejs, save this file locally, and run with:
//    node OklabExperiments.js
// No other dependencies are required to use this.
// Thanks to some helpful folks in the generative art community for helping me better understand what's happening with this.

// My toddler smacked the keyboard with a piece of cardboard and made me accidentally type:
// zaser~

// Some color / math code tweaked from functions found in this repo: https://github.com/mattdesl/tiny-artblocks

function clamp(value, min, max) {
  return Math.max(Math.min(value, max), min);
}

// returns an object literal which is r, g and b integers from an RGB hex string:
function hexToRGB(str) {
  var hex = str.replace("#", "");
  // NOTE: This can be removed for brevity if you stick with 6-character codes:
  // if (hex.length === 3) {hex = hex[0] + hex[0] + hex[1] + hex[1] + hex[2] + hex[2];}
  var num = parseInt(hex, 16);
  return { r: num >> 16, g: (num >> 8) & 255, b: num & 255 };
}

// converts RGB integer values to hex:
function rgbToHex({ r, g, b }) {
  // NEXT LINE OPTIONAL, depending on your purposes:
  // var r = clamp(r, 0, 255); var g = clamp(g, 0, 255); var b = clamp(b, 0, 255);
  return "#" + (b | (g << 8) | (r << 16) | (1 << 24)).toString(16).slice(1);
}

// correlary of first psuedocode block here (f_inv) : https://bottosson.github.io/posts/colorwrong/#what-can-we-do%3F ; "applying the inverse of the sRGB nonlinear transform function.." -- keeping the abbreviated syntax of arrow functions and ? : if/then, despite that they confuse and stretch my noob brain:
const gammaToLinear = (c) =>
  c >= 0.04045 ? Math.pow((c + 0.055) / 1.055, 2.4) : c / 12.92;
// correlary of the first " : "..then switching back" :
const linearToGamma = (c) =>
  c >= 0.0031308 ? 1.055 * Math.pow(c, 1 / 2.4) - 0.055 : 12.92 * c;

// Lab coordinates (parameters):
// L = Lightness (0 (black) to ?? (diffuse white)
// a = green (at negative -- is there a lower bound?) to red (positive)
// b = blue (at negative) to yellow (at positive).
// You can manually construct an object literal to pass to this function this way:
// labVals = {L: 0.75, a: 0.7, b: 0.2};
// sRGBresultObjectLiteral = oklabToSRGB(labVals);
// You can also construct that as just {0.75, 0.7, 0.2}, and still pass it and it will work, apparently
// "..Oklab is represented as an object {L, a, b} where L is between 0 and 1 for normal SRGB colors. a and b have a less clearly defined range, but will normally be between -0.5 and +0.5. Neutral gray colors will have a and b at zero (or very close)." re: https://www.npmjs.com/package/oklab
// numbers updated from C++ example at https://bottosson.github.io/posts/oklab/ as updated 2021-01-25
// helpful references:
// https://observablehq.com/@sebastien/srgb-rgb-gamma
// Other references: https://matt77hias.github.io/blog/2018/07/01/linear-gamma-and-sRGB-color-spaces.html
// Takeaway: before manipulating color for sRGB (gamma-corrected or gamma compressed), convert it to linear RGB or a linear color space. Then do the manipulation, then convert it back to sRGB.
function rgbToOklab({ r, g, b }) {
  // This is my undersanding: JavaScript canvas and many other virtual and literal devices use gamma-corrected (non-linear lightness) RGB, or sRGB. To convert sRGB values for manipulation in the Oklab color space, you must first convert them to linear RGB. Where Oklab interfaces with RGB it expects and returns linear RGB values. This next step converts (via a function) sRGB to linear RGB for Oklab to use:
  r = gammaToLinear(r / 255);
  g = gammaToLinear(g / 255);
  b = gammaToLinear(b / 255);
  // This is the Oklab math:
  var l = 0.4122214708 * r + 0.5363325363 * g + 0.0514459929 * b;
  var m = 0.2119034982 * r + 0.6806995451 * g + 0.1073969566 * b;
  var s = 0.0883024619 * r + 0.2817188376 * g + 0.6299787005 * b;
  // Math.crb (cube root) here is the equivalent of the C++ cbrtf function here: https://bottosson.github.io/posts/oklab/#converting-from-linear-srgb-to-oklab
  l = Math.cbrt(l);
  m = Math.cbrt(m);
  s = Math.cbrt(s);
  return {
    L: l * +0.2104542553 + m * +0.793617785 + s * -0.0040720468,
    a: l * +1.9779984951 + m * -2.428592205 + s * +0.4505937099,
    b: l * +0.0259040371 + m * +0.7827717662 + s * -0.808675766,
  };
}

function oklabToSRGB({ L, a, b }) {
  var l = L + a * +0.3963377774 + b * +0.2158037573;
  var m = L + a * -0.1055613458 + b * -0.0638541728;
  var s = L + a * -0.0894841775 + b * -1.291485548;
  // The ** operator here cubes; same as l_*l_*l_ in the C++ example:
  l = l ** 3;
  m = m ** 3;
  s = s ** 3;
  var r = l * +4.0767416621 + m * -3.3077115913 + s * +0.2309699292;
  var g = l * -1.2684380046 + m * +2.6097574011 + s * -0.3413193965;
  var b = l * -0.0041960863 + m * -0.7034186147 + s * +1.707614701;
  // Convert linear RGB values returned from oklab math to sRGB for our use before returning them:
  r = 255 * linearToGamma(r);
  g = 255 * linearToGamma(g);
  b = 255 * linearToGamma(b);
  // OPTION: clamp r g and b values to the range 0-255; but if you use the values immediately to draw, JavaScript clamps them on use:
  r = clamp(r, 0, 255);
  g = clamp(g, 0, 255);
  b = clamp(b, 0, 255);
  // OPTION: round the values. May not be necessary if you use them immediately for rendering in JavaScript, as JavaScript (also) discards decimals on render:
  r = Math.round(r);
  g = Math.round(g);
  b = Math.round(b);
  return { r, g, b };
}

// A DEV TEST that determined min and max range for oklab L, a, and b as found via conversion from all possible RGB256 values (I want those ranges for the curiosity of generating random colors in Oklab space but converting the results to rgb):
// var counter = 0;
// var aMin = 0, aMax = 0, bMin = 0, bMax = 0, lMin = 50; lMax = -50;
// for (var r = 0; r < 256; r++) {
//   for (var g = 0; g < 256; g++) {
//     for (var b = 0; b < 256; b++) {
//       counter += 1;
//       var sRGBcolor = {r, g, b};
//       var labColor = rgbToOklab(sRGBcolor);
//       if (labColor.L < lMin) {lMin = labColor.L;}
//       if (labColor.L > lMax) {lMax = labColor.L;}
//       if (labColor.a < aMin) {aMin = labColor.a;}
//       if (labColor.a > aMax) {aMax = labColor.a;}
//       if (labColor.b < bMin) {bMin = labColor.b;}
//       if (labColor.b > bMax) {bMax = labColor.b;}
//     }
//   }
// }
// console.log("counter:", counter);
// console.log("lMin:", lMin, "lMax:", lMax, "aMin:", aMin, "aMax:", aMax, "bMin:", bMin, "bMax:", bMax);
// results without and with gammaToLinear -> linearToGamma and other "OPTIONAL" calculations:
// without extra calculations       with extra calculations
// lMin: 0                          0
// lMax: 254.99999832226274         0.9999999934735462
// aMin: -118.66066044645729        -0.23388757418790818
// aMax: 118.66066050838631         0.27621639742350523
// bMin: -108.13237868927439        -0.3115281476783751
// bMax: 108.1323881808137          0.19856975465179516

// A DEV TEST that prints so many round trips from sRGB -> Oklab -> sRGB to verify that randomly generated colors remain the same after round trip. They do.
// function getRNDrgbColor() {
//   const rndRGB = n => Math.round(Math.random() * 256);
//   return {r: rndRGB(), g: rndRGB(), b: rndRGB()}
// }
// for (let i = 0; i < 42; i++) {
//   rndRGBvals = getRNDrgbColor();
//   console.log(rndRGBvals, "->");
//   labVals = rgbToOklab(rndRGBvals);    // can pass rndRGBvals bcse it is like {r: val, g: val, b: val}
//   console.log(oklabToSRGB(labVals), "<-\n");
// }

// POSSIBILITY FOR GRAYSCALE CALC; you can use the value returned from this with a = 0, b = 0 (for oklab) to convert RGB to gray:
// function linearrgbToOklabLightness({r, g, b}) {
//   const l = Math.cbrt(r * +0.4122214708 + g * +0.5363325363 + b * +0.0514459929);
//   const m = Math.cbrt(r * +0.2119034982 + g * +0.6806995451 + b * +0.1073969566);
//   const s = Math.cbrt(r * +0.0883024619 + g * +0.2817188376 + b * +0.6299787005);
//
//   return l * +0.2104542553 + m * +0.7936177850 + s * -0.0040720468;
// };

// returns an array which is a range of values at N intervals over range min-max (interpolations), inclusive of max:
function getLerpRange(min, max, N) {
  dividend = (max - min) / (N - 1);
  var arr = [];
  currentVal = min;
  for (var i = 0; i < N; i++) {
    arr.push(currentVal);
    currentVal += dividend;
  }
  return arr;
}

// returns an array of sRGB hex color codes, from interpolation in numberOfColors steps through Oklab space (and converted back to sRGB), including start and end value:
function getRGBlerpRangeInOklabSpace(
  startRGBhexColor,
  endRGBhexColor,
  numberOfColors,
) {
  var lasInterpolaciones = [];
  var startLabVals = rgbToOklab(hexToRGB(startRGBhexColor));
  var endLabVals = rgbToOklab(hexToRGB(endRGBhexColor));
  var lValsArr = getLerpRange(startLabVals.L, endLabVals.L, numberOfColors);
  var aValsArr = getLerpRange(startLabVals.a, endLabVals.a, numberOfColors);
  var bValsArr = getLerpRange(startLabVals.b, endLabVals.b, numberOfColors);
  for (var i = 0; i < numberOfColors; i++) {
    // I could nest function calls here for brevity . . .
    var tmpOklabVals = { L: lValsArr[i], a: aValsArr[i], b: bValsArr[i] };
    var tmpsRGBvals = oklabToSRGB(tmpOklabVals);
    // var hexVals = rgbToHex(tmpsRGBvals);
    lasInterpolaciones.push(tmpsRGBvals);
  }
  return lasInterpolaciones;
}

testInterpolaciones = getRGBlerpRangeInOklabSpace("#0000ff", "#0f0f0ff00", 7);
for (var element of testInterpolaciones) {
  //console.log(element);
}
// compare the middle value of that to the blue -> yellow gradient through Oklab space here: https://raphlinus.github.io/color/2021/01/18/oklab-critique.html -- they match. This (and other tests) verifies this is doing things correctly from sRGB -> linear RGB -> Oklab -> transformation -> linear RGB -> sRGB.

let picoPrimary = { L: 0.529, a: 0.1254, b: 242.05 };

const prompt = require("prompt-sync")();
let colourString = "";
let color = { L: 0, a: 0, b: 0 };

while (true) {
  console.log("colour: ");
  colourString = prompt("");
  console.clear();

  let color = { L: 0, a: 0, b: 0 };
  let opacity = 1;

  if (colourString.includes("rgba")) {
    let rgblist = colourString
      .match(
        /^(0|[1-9]\d*)(\.\d+)?$, ^(0|[1-9]\d*)(\.\d+)?$, ^(0|[1-9]\d*)(\.\d+)?$, ^(0|[1-9]\d*)(\.\d+)?$/,
      )
      .toString()
      .split(",")
      .map((e) => Number(e));
    color = rgbToOklab({ r: rgblist[0], g: rgblist[1], b: rgblist[2] });
    opacity = rgblist[3];
  } else if (colourString.includes("rgb")) {
    let rgblist = colourString
      .match(/( (\d+)?(\.\d+) ), ( (\d+)?(\.\d+) ), ((\d+)?(\.\d+))/)
      .toString()
      .split(",")
      .map((e) => Number(e));

    console.log(rgblist);
    color = rgbToOklab({ r: rgblist[0], g: rgblist[1], b: rgblist[2] });
  } else if (colourString.includes("#")) {
    let hex = colourString.match(/#....../).toString();
    color = rgbToOklab(hexToRGB(hex));
  }

  let { Factor, Color } = getColorToMixAndFactor(picoPrimary, color);

  let colourToPrint = "";

  if (Factor >= 1) {
    colourToPrint = `oklab(${color.L}, ${color.a}, ${color.b})`;
  } else if (opacity != 1) {
    colourToPrint = `color-mix(in oklab, color-mix(in oklab, var(--pico-primary), oklab(${Color.L}, ${Color.a}, ${Color.b} ${Factor * 100}%), transparent ${opacity * 100}%)`;
  } else if (opacity == 0) {
    colourToPrint = "transparent";
  } else {
    colourToPrint = `color-mix(in oklab, var(--pico-primary), oklab(${Color.L}, ${Color.a}, ${Color.b} ${Factor * 100}%)`;
  }

  console.log(colourToPrint);
}
