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
      files: ['web/dist']
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
              'web/src/js/jquery-3.2.1.min.js',
              'web/src/js/lodash.min.js',
              'web/src/js/moment.min.js',
              'web/src/js/fullcalendar.min.js',
              'web/src/js/vex.combined.min.js',
              'web/src/js/app.js',
              'web/src/js/**.js'
          ],
          dest: 'web/dist/js/<%= pkg.name %>.js'
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
          'web/dist/css/<%= pkg.name %>.min.css': ['web/src/css/*.css']
        }
      }
    },
    insert:
      {
        options: {},
         js: {
                 src:  '<%= uglify.dist.dest %>',
                 dest:   '<%= concat.disthtml.dest %>',
                 match: "<!-- grunt-insert:js -->"
         },
         css: {
                 src:  'web/dist/css/<%= pkg.name %>.min.css',
                 dest:   '<%= concat.disthtml.dest %>',
                 match: "<!-- grunt-insert:css -->"
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
  grunt.loadNpmTasks('grunt-insert');

  // Default task.
  grunt.registerTask('default', ['clean', 'concat', 'uglify', 'cssmin', 'insert']);

};
