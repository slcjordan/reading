

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
            if ((end - start) < wait){
                return false;
            }
            return true;
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
        return vex.dialog.open(_.merge(options, {
            beforeClose: debounce(500),
            focusFirstInput: false,
            input: _.map(options.input, renderInput).join(''),
            buttons: buttons
        }));
    }

    function dialogPlugin(vex) {
        return {
            name: 'plan',
            open: _.partial(dialog, vex)
        };
    }

    vex.registerPlugin(dialogPlugin);

    function populate(target, start, sessions, options){
        var dates = _.map(
            _.times(sessions.length),
            function(incr){
                return {'start': start.clone().add(incr, 'day')};
            }
        );
        var events = _.zipWith(sessions, dates, _.merge);
        _.forEach(events, _.partial(_.set, _, 'events', events));
        $(target).fullCalendar('renderEvents', events, true);  // stick
        $(target).fullCalendar('unselect');
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
                $.getJSON('plan', options, function(sessions){
                    populate($('#calendar'), start, _.defaultTo(sessions, []), options);
                });
            }
        });
    }

    function eventUpdateFlow(event){
        vex.plan.open({
            message: event.title + ' on ' + event.start.format('ddd, MMM D'),
            input: [
            ],
            no:  'OK'
        });
    }

    $('#calendar').fullCalendar({
        editable: true,
        header: {
          left: 'prev,next',
            center: 'title',
            right: 'month,agendaWeek'
        },
        selectable: false,
        selectHelper: true,
        dayClick: planCreateFlow,
        eventClick: eventUpdateFlow
    });
});
