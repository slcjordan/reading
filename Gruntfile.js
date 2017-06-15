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
      distjs: {
        src: [
            'web/src/js/jquery-3.2.1.min.js',
            'web/src/js/lodash.min.js',
            'web/src/js/moment.min.js',
            'web/src/js/fullcalendar.min.js',
            'web/src/js/vex.combined.min.js',
            'web/src/js/app.js'
        ],
        dest: 'web/dist/js/<%= pkg.name %>.js'
      },
      disthtml: {
        src: ['web/src/index.html'],
        dest: 'web/dist/index.html'
      },
    },
    uglify: {
      options: {
        banner: '<%= banner %>'
      },
      dist: {
        src: '<%= concat.distjs.dest %>',
        dest: 'web/dist/js/<%= pkg.name %>.min.js'
      },
    },
    cssmin: {
      options: {
        mergIntoShorthands: false,
        roundingPrecision: -1
      },
      target: {
        files: {
          'web/dist/css/<%= pkg.name %>.min.css': ['web/src/css/*.css']
        },
      },
    },
    jshint: {
      options: {
        jshintrc: true,
        reporterOutput: ""
      },
      gruntfile: {
        src: 'Gruntfile.js'
      },
      src: {
        src: ['web/src/src/app.js']
      },
      test: {
        src: ['web/test/**/*.js']
      },
    },
    watch: {
      gruntfile: {
        files: '<%= jshint.gruntfile.src %>',
        tasks: ['jshint:gruntfile']
      },
      src: {
        files: '<%= jshint.src.src %>',
        tasks: ['jshint:src']
      },
      test: {
        files: '<%= jshint.test.src %>',
        tasks: ['jshint:test']
      },
    },
  });

  // These plugins provide necessary tasks.
  grunt.loadNpmTasks('grunt-contrib-clean');
  grunt.loadNpmTasks('grunt-contrib-concat');
  grunt.loadNpmTasks('grunt-contrib-uglify');
  grunt.loadNpmTasks('grunt-contrib-jshint');
  grunt.loadNpmTasks('grunt-contrib-watch');
  grunt.loadNpmTasks('grunt-contrib-cssmin');

  // Default task.
  grunt.registerTask('default', ['jshint', 'clean', 'concat', 'uglify', 'cssmin']);

};
