export const format = {
	/**
	 * @param {number} num,
	 * @param {boolean?} add_plus
	 * @param {number?} decimal_places
	 * @returns {string}
	 */
	number: function (num, add_plus, decimal_places) {
		if (decimal_places == undefined) {
			decimal_places == 2;
		}

		if (!add_plus || num <= 0) {
			return num.toLocaleString('de-CH', {
				maximumFractionDigits: decimal_places,
				minimumFractionDigits: decimal_places
			});
		}

		return `+${num.toLocaleString('de-CH', { maximumFractionDigits: decimal_places, minimumFractionDigits: decimal_places })}`;
	},

	/**
	 * @param {number | undefined} num,
	 * @param {boolean?} add_plus
	 * @param {number?} decimal_places
	 * @returns {string}
	 */
	currency: function (num, add_plus, decimal_places) {
		if (num == undefined) {
			return '';
		}
		if (decimal_places == undefined) {
			/** @type {number} */
			decimal_places == 0;
		}
		if (add_plus === undefined) {
			/** @type {boolean} */
			add_plus = false;
		}

		return `${this.number(num, add_plus, decimal_places)} CHF`;
	},

	/**
	 * @param {string} string
	 * @returns {string}
	 */
	capitaliseFirstLetter: function (string) {
		return string[0].toUpperCase() + string.slice(1);
	},
	/**
	 * @param {string} string
	 * @returns {string}
	 */
	titleCase: function (string) {
		let titleString = '';
		for (let word of string.split(' ')) {
			titleString += this.capitaliseFirstLetter(word);
			titleString += ' ';
		}

		return titleString.slice(0, titleString.length - 1);
	}
};
