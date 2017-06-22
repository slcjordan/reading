

$(document).ready(function() {

    function includeTemplate(a, b){
        return function(options){
            var include = b(options);
            return a(_.set(options, 'include', include));
        };
    }

    function optionList(data){
        return _
            .chain(data.options)
            .map(function(val){ return {'value': val};})
            .map(_.template('<option value="<%= value %>"><% print(_.startCase(value)) %></option>'))
            .value()
            .join('');
    }
        
    var label = _.template(
        [   
            '<div class="vex-custom-field-wrapper">',
                '<label for="<%= name %>"><% print(_.startCase(name)) %>: </label>',
                '<%= include %>',
            '</div>'
        ].join('')
    );

    var select = includeTemplate(
        _.template('<select name="<%= name %>" required><%= include %></select>'),
        optionList
    );

    var number = _.template('<input name="<%= name %>" min="1" type="number" value="<%= initial %>" required />');

    function debounce(wait) {
        var start = new Date().getTime();
        return function(){
            var end = new Date().getTime();
            return end - start > wait;
        };
    }

    function dialog(vex, options){
        var buttons = [];
        if (_.has(options, 'yes')){
            buttons.push(_.merge({}, vex.dialog.buttons.YES, { text: options.yes }));
        }
        if (_.has(options, 'no')){
            buttons.push(_.merge({}, vex.dialog.buttons.NO, { text: options.no }));
        }
        var renderInput = function(input){
            if (_.has(input, 'options')){
                return includeTemplate(label, select)(input);
            } else if (_.isNumber(input.initial)) {
                return includeTemplate(label, number)(input);
            }
            return '';
        };
        var cb = options.callback;
        return vex.dialog.open(_.merge(options, {
            beforeClose: debounce(500),
            focusFirstInput: false,
            input: _.map(options.input, renderInput).join(''),
            buttons: buttons,
            callback: function(data){
                if (data===false){
                    return;
                }
                cb(data);
            }
        }));
    }

    function dialogPlugin(vex) {
        return {
            name: 'plan',
            open: _.partial(dialog, vex)
        };
    }

    vex.registerPlugin(dialogPlugin);

    function sort(dates){
        var unix = function(e){return e.start.unix();}
        var sorted = _.sortedUniqBy(_.sortBy(dates, unix), unix);
        return _.each(sorted, _.partial(_.set, _, 'index'));
    }

    function makeDates(start, days){
        if (days <= 0) {
            return [];
        }
        return _.map(
            _.times(days),
            function(incr){
                return {
                    'start': start.clone().add(incr, 'day')
                };
            }
        );
    }

    function AJAX(url, start, data){
        var reading = {dates: sort(makeDates(start, data.days))};
        return {
            url: url,
            error: function(){console.log('could not fetch the reading plan.');},
            success: function(events){return _.zipWith(reading.dates, events, _.merge);},
            reading: reading,
            data: data
        };
    }

    function planCreateFlow(start, end){
        vex.plan.open({
            message: 'Starting ' + start.format('ll'),
            input: [
                {
                    name: 'book',
                    options: [
                        'book-of-mormon',
                        'new-testament',
                        'old-testament',
                        'doctrine-and-covenants',
                        'pearl-of-great-price'
                    ]
                },
                {
                    name: 'days',
                    initial: 90
                },
                {
                    name: 'breakdown',
                    options: [
                        'chapter',
                        'verse'
                    ]
                }
            ],
            yes: 'Create',
            no:  'Cancel',
            callback: function(options){
                // TODO Dbug
                $('#calendar').fullCalendar('addEventSource', AJAX('plan', start, options));
                $('#calendar').fullCalendar('unselect');
            }
        });
    }

    function remove(event){
        event.source.data.days--;
        event.source.reading.dates = sort(_.filter(
            event.source.reading.dates,
            function(d){
                return event.index !== d.index;
            }
        ));
    }

    function eventUpdateFlow(event){
        vex.plan.open({
            message: event.title + ' on ' + event.start.format('ddd, MMM D'),
            input: [],
            yes: 'Delete',
            no:  'Cancel',
            callback: function(){
                remove(event);
                $('#calendar').fullCalendar('refetchEventSources', event.source);
            }
        });
    }

    $('#calendar').fullCalendar({
        editable: true,
        header: {
          left: 'prev,next',
            center: 'title',
            right: ''
        },
        selectable: false,
        selectHelper: true,
        eventOverlap: false,
        dayClick: planCreateFlow,
        eventClick: eventUpdateFlow,
        eventDrop: function(event){
            remove(event);
            var dates = makeDates(event.start, 1);
            event.source.reading.dates = sort(_.concat(event.source.reading.dates, dates));
            event.source.data.days = _.size(event.source.reading.dates);
            $('#calendar').fullCalendar('refetchEventSources', event.source);
        },
        eventResize: function(event){
            remove(event);
            var days = event.end.diff(event.start, 'days');
            var dates = makeDates(event.start, days);
            event.source.reading.dates = sort(_.concat(event.source.reading.dates, dates));
            event.source.data.days = _.size(event.source.reading.dates);
            $('#calendar').fullCalendar('refetchEventSources', event.source);
        }
    });
});
vex.defaultOptions.className = 'vex-theme-wireframe';
