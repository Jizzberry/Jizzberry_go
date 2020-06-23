var videosHound = new Bloodhound({
    datumTokenizer: Bloodhound.tokenizers.whitespace,
    queryTokenizer: Bloodhound.tokenizers.whitespace,
    remote: {
        url: "/api/files?file_name=%QUERY",
        wildcard: '%QUERY',
        filter: function (response) {
            console.log(response);
            return response;
        }
    }
});

var actorsHound = new Bloodhound({
    datumTokenizer: Bloodhound.tokenizers.whitespace,
    queryTokenizer: Bloodhound.tokenizers.whitespace,
    remote: {
        url: "/api/actor_details?name=%QUERY",
        wildcard: '%QUERY',
        filter: function (response) {
            console.log(response);
            return response;
        }
    }
});

videosHound.initialize();
actorsHound.initialize();

$('#bloodhound .typeahead').typeahead(
    {
        hint: true,
        minLength: 0,
        highlight: true,
    },
    {
        name: 'files',
        displayKey: 'file_name',
        source: videosHound.ttAdapter(),
        templates: {
            header: '<h6 class="group-name">Scenes</h6>',
            suggestion: Handlebars.compile('<p><a href="/Jizzberry/scene/{{generated_id}}" style="display: block;">{{file_name}}</a></p>')
        },
        limit: 20
    },
    {
        name: 'actors',
        displayKey: 'name',
        source: actorsHound.ttAdapter(),
        templates: {
            header: '<h6 class="group-name">Actors</h6>',
            suggestion: Handlebars.compile('<p><a href="/Jizzberry/actors/{{actor_id}}" style="display: block;">{{name}}</a></p>')
        },
        limit: 20
    }
);