
$(document).ready(function() {

    $('#calendar').fullCalendar({
        editable: true,
        // header
        header: {
          left: 'prev,next today',
            center: 'title',
            right: 'month,agendaWeek'
        },
        selectable: true,
        selectHelper: true,
        select: function(start, end) {
            vex.dialog.open({
                message: 'Starting ' + start.format('ll'),
                input: [
                    '<select name="book" required>',
                         '<option value="book-of-mormon">Book of Mormon</option>',
                         '<option value="new-testament">New Testament</option>',
                         '<option value="old-testament">Old Testament</option>',
                         // '<option value="doctrine-and-covenants">Doctrine and Covenants</option>',
                         // '<option value="pearl-of-great-price">Pearl of Great Price</option>',
                     '</select>',
                    '<div class="vex-custom-field-wrapper">',
                        '<label for="days">days: </label>',
                        '<input name="days" min="1" max="800" type="number" placeholder="Number of Days" value="90" required />',
                    '</div>',
                    '<div class="vex-custom-field-wrapper">',
                        '<label for="breakdown">broken down by: </label>',
                        '<select name="breakdown">',
                            '<option value="chapter">Chapter/Section</option>',
                            '<option value="verse">Verse</option>',
                        '</select>',
                    '</div>'
                ].join(''),
                buttons: [
                    $.extend({}, vex.dialog.buttons.YES, { text: 'Create' }),
                    $.extend({}, vex.dialog.buttons.NO, { text: 'Cancel' })
                ],
                callback: function (data) {
                    console.log(data);
                    var eventData;
                    $.ajax(
                        'http://localhost:8080/plan',
                        {
                        data: data,
                        success: function(result, status, jsqxhr){
                            var events;
                            events = JSON.parse(result).map(function(v, i){
                                v.start = start.clone().add(i, 'day');
                                v.title = v.Title;
                                return v;
                            });
                            console.log(events);
                            $('#calendar').fullCalendar(
                                'renderEvents',
                                events,
                                true
                            );
                        },
                        error: function(jqxhr, status, err){
                            console.log(status);
                            console.log(err); }
                        });

                    $('#calendar').fullCalendar('unselect');
                }
            });
        }
    });

});
