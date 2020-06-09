	$(document).ready(function () {
				// Defining the local dataset
				var pornstar = [
					"Mia Khalifa",
					"Lana Rhoades",
					"Leah Gotti",
					"Abella Danger",
					"Jessa Rhoades",
					"Dani Daniels",
					"Mia Melano ",
					"Natasha",
					"Riley",
				];

				// Constructing the suggestion engine
				var pornstar = new Bloodhound({
					datumTokenizer: Bloodhound.tokenizers.whitespace,
					queryTokenizer: Bloodhound.tokenizers.whitespace,
					local: pornstar,
				});

				// Initializing the typeahead
				$(".typeahead").typeahead(
					{
						hint: true,
						highlight: true /* Enable substring highlighting */,
						minLength: 1 /* Specify minimum characters required for showing result */,
					},
					{
						name: "pornstar",
						source: pornstar,
					}
				);
			});