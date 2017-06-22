'use strict';

module.exports = function(grunt) {

  // Project configuration.
  grunt.initConfig({
    // Metadata.
    pkg: grunt.file.readJSON('web/reading.jquery.json'),
    banner: '/*! <%= pkg.title || pkg.name %> - v<%= pkg.version %> - ' +
      '<%= grunt.template.today("yyyy-mm-dd") %>\n' +
      '<%= pkg.homepage ? "* " + pkg.homepage + "\\n" : "" %>' +
      '* Copyright (c) <%= grunt.template.today("yyyy") %> <%= pkg.author.name %>;' +
      ' Licensed <%= _.pluck(pkg.licenses, "type").join(", ") %> */\n',
    // Task configuration.
    clean: {
      files: ['web/dist', 'web/dev']
    },
    concat: {
      options: {
        stripBanners: true,
        separator: ';\n'
      },
      disthtml: {
        src: [
            'web/src/index.html'
        ],
        dest: 'web/dist/index.html'
      },
      distjs: {
          src: [
              'web/src/js/jquery-3.2.1.js',
              'web/src/js/lodash.js',
              'web/src/js/moment.js',
              'web/src/js/fullcalendar.js',
              'web/src/js/vex.combined.js',
              'web/src/js/app.js',
              'web/src/js/analytics.js'
          ],
          dest: 'web/dist/js/<%= pkg.name %>.js'
      }
    },
    copy:{
      dev: {
        files: [
          {expand: true, cwd: 'web/src', src: ['**'], dest: 'web/dev/', filter: 'isFile'}
        ]
      }
    },
    uglify: {
      options: {
        banner: '<%= banner %>'
      },
      dist: {
        src: '<%= concat.distjs.dest %>',
        dest: 'web/dist/js/<%= pkg.name %>.min.js'
      }
    },
    cssmin: {
      options: {
        mergIntoShorthands: false,
        roundingPrecision: -1
      },
      target: {
        files: {
          'web/dist/css/<%= pkg.name %>.min.css': [
              'web/src/css/fullcalendar.css',
              'web/src/css/vex.css',
              'web/src/css/vex-theme-wireframe.css',
              'web/src/css/app.css'
          ]
        }
      }
    },
    writefile: {
        options: {
            paths: {
                js: {
                    cwd: 'web/src',
                    src: [
                          'js/jquery-3.2.1.js',
                          'js/lodash.js',
                          'js/moment.js',
                          'js/fullcalendar.js',
                          'js/vex.combined.js',
                          'js/app.js',
                          'js/app.js'
                    ]
                },
                css: {
                    cwd: 'web/src',
                    src: [
                        'css/fullcalendar.css',
                        'css/vex.css',
                        'css/vex-theme-wireframe.css',
                        'css/app.css'
                    ]
                }
            }
        },
        devjs: {
            src: 'web/src/templates/js.hbs',
            dest: 'web/dev/js.html'
        },
        devcss: {
            src: 'web/src/templates/css.hbs',
            dest: 'web/dev/css.html'
        }
    },
    insert:
      {
        options: {},
         distjs: {
                 src:  '<%= uglify.dist.dest %>',
                 dest:   '<%= concat.disthtml.dest %>',
                 match: '<!-- grunt-insert:js -->'
         },
         devjs: {
                 src:  '<%= writefile.devjs.dest %>',
                 dest:   'web/dev/index.html',
                 match: '<script type="text/javascript"><!-- grunt-insert:js --></script>'
         },
         distcss: {
                 src:  'web/dist/css/<%= pkg.name %>.min.css',
                 dest:   '<%= concat.disthtml.dest %>',
                 match: '<!-- grunt-insert:css -->'
         },
         devcss: {
                 src:  '<%= writefile.devcss.dest %>',
                 dest:   'web/dev/index.html',
                 match: '<style><!-- grunt-insert:css --></style>'
         }
    },
    watch: {
      src: {
        files: '<%= jshint.src.src %>',
        tasks: ['jshint:src']
      }
    }
  });

  // These plugins provide necessary tasks.
  grunt.loadNpmTasks('grunt-contrib-clean');
  grunt.loadNpmTasks('grunt-contrib-concat');
  grunt.loadNpmTasks('grunt-contrib-uglify');
  grunt.loadNpmTasks('grunt-contrib-jshint');
  grunt.loadNpmTasks('grunt-contrib-watch');
  grunt.loadNpmTasks('grunt-contrib-cssmin');
  grunt.loadNpmTasks('grunt-writefile');
  grunt.loadNpmTasks('grunt-insert');
  grunt.loadNpmTasks('grunt-contrib-copy');

  // Default task.
  grunt.registerTask('default', ['clean', 'copy', 'concat', 'uglify', 'cssmin', 'writefile', 'insert']);

};
