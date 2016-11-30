var _ = require('lodash');
var gulp = require('gulp');
var less = require('gulp-less');
var gutil = require('gulp-util');
var concat = require('gulp-concat');
var uglify = require('gulp-uglify');
var cssnano = require('gulp-cssnano');
var replace = require('gulp-replace');
var rename = require('gulp-rename');
var path = require('path');
var webpack = require('webpack');
var webpackConfig = require('./webpack.config.js');
var exec = require('child_process').exec;
var walk = require('walk');
var fs = require('fs');
var request = require('request');
var async = require('async');
var browserSync = require('browser-sync').create();

gulp.task('assets', function () {
  gulp
    .src('./assets/**/*.*')
    .pipe(gulp.dest('build/statics'));
});

gulp.task('translation-push', function (done) {
  async.auto({
    xgettext: function (next) {
      exec([
        'xgettext',
        '--force-po',
        '-o',
        'po/putio.pot',
        '--language="Javascript"',
        '--keyword="_t"',
        '--keyword="_tn:1,2"',
        '--keyword="_tc:1,2c"',
        '--keyword="_tnc:1,2,4c"',
        '--from-code="utf-8"',
        'src/**/*.js'
      ].join(' '), next);
    },
    getPot: function (next) {
      request({
        method: 'GET',
        uri: [
          transifex.url,
          transifex.project,
          'resource',
          transifex.resource,
          'content',
          '?file'
        ].join('/'),
        auth: transifex.auth,
      }, function (err, response, body) {
        if (err) {
          return next(err);
        }

        fs.writeFile(__dirname + '/po/putio.pot.orig', body, next)
      })
    },
    dirty: ['getPot', 'xgettext', function (results, next) {
      exec([
        'msgcomm',
        '-u',
        '--omit-header',
        '"' + __dirname + '/po/putio.pot.orig' + '"',
        '"' + __dirname + '/po/putio.pot' + '"',
      ].join(' '), function (err, stdout) {
        if (err) {
          return next(err)
        }

        return next(null, !!stdout)
      })
    }],
    push: ['dirty', function (results, next) {
      if (!results.dirty) {
        return next()
      }

      request({
        method: 'PUT',
        uri: [
          transifex.url,
          transifex.project,
          'resource',
          transifex.resource,
          'content'
        ].join('/'),
        auth: transifex.auth,
        formData: {
          content: fs.createReadStream(__dirname + '/po/putio.pot')
        },
      }, function (err, response, body) {
        if (err) {
          return next(err)
        }

        return next();
      })
    }]
  }, done)
});

gulp.task('translation-pull', ['translation-push'], function (done) {
  async.auto({
    languages: function (next) {
      request({
        method: 'GET',
        uri: transifex.url + '/' + transifex.project + '/languages/',
        uri: [
          transifex.url,
          transifex.project,
          'languages',
        ].join('/'),
        json: true,
        auth: transifex.auth,
      }, function (err, response, body) {
        if (err) {
          return next(err);
        }

        return next(null, body);
      });
    },
    translations: ['languages', function (results, next) {
      async.mapSeries(results.languages, function (lang, next) {
        request({
          method: 'GET',
          uri: [
            transifex.url,
            transifex.project,
            'resource',
            transifex.resource,
            'translation',
            lang.language_code,
            '?file'
          ].join('/'),
          auth: transifex.auth,
        }, function (err, response, body) {
          if (err) {
            return next(err);
          }

          return next(null, {
            code: lang.language_code,
            data: body
          });
        })
      }, next)
    }],
    save: ['translations', function (results, next) {
      async.eachSeries(results.translations, function (translation, next) {
        fs.writeFile(
          __dirname + '/po/' + translation.code + '.po',
          translation.data,
          next
        )
      }, next)
    }],
    tojson: ['save', function (results, next) {
      var walker  = walk.walk("./po", {
        followLinks: false
      })

      walker.on("file", (root, stat, next) => {
        var srcPath = path.resolve(root, stat.name)
        var ext = path.extname(srcPath)
        var dstPath = "assets/locale/"

        if (ext == '.pot') {
          // putio.pot is the source language, hence en.json
          dstPath += 'en.json'
        } else {
          dstPath += path.basename(stat.name, ext) + '.json'
        }

        exec([
          'node_modules/po2json/bin/po2json',
          "'" + srcPath + "'",
          "'" + dstPath + "'",
          '-f',
          'jed1.x',
          '-d',
          'putio',
        ].join(' '), next);
      });

      walker.on("end", next);
    }]
  }, done);
});

gulp.task('less', function () {
  gulp.src('./src/app/style.less')
    .pipe(less())
    .pipe(rename({ suffix: '.min' }))
    .pipe(cssnano())
    .pipe(gulp.dest('./build/statics/css'));
});

gulp.task('less-dev', function () {
  gulp.src('./src/app/style.less')
    .pipe(less())
    .pipe(gulp.dest('./build/statics/css'))
    .pipe(browserSync.stream())
});

gulp.task('html-dev', function () {
  gulp
    .src('index.html')
    .pipe(gulp.dest('build'))
});

gulp.task('html', function () {
  gulp
    .src('index.html')
    .pipe(replace(/index.js/g, 'index.min.js'))
    .pipe(replace(/style.css/g, 'style.min.css'))
    .pipe(replace(/development/g, 'production'))
    .pipe(gulp.dest('build'));
})

gulp.task('build', [
  'less',
  'html',
  'assets',
  'libjs',
], function (next) {
  var webpackProdConfig = Object.create(webpackConfig);
  webpackProdConfig.devtool = 'sourcemap';
  webpackProdConfig.plugins = webpackProdConfig.plugins || []
  webpackProdConfig.plugins.push(
    new webpack.DefinePlugin({
      'process.env.NODE_ENV': JSON.stringify('production'),
    })
  )
  webpackProdConfig.plugins.push(
    new webpack.optimize.UglifyJsPlugin()
  )
  webpackProdConfig.output.filename = 'index.min.js';
  var wp = webpack(webpackProdConfig).run(next);
});

gulp.task('html-watch', ['html-dev'], function (done) {
  browserSync.reload()
  done()
})

gulp.task('build-dev', [
  'webpack:build-dev',
  'less-dev',
  'html-dev',
  'assets',
  'libjs-dev',
], function () {
  browserSync.init({
    ui: false,
    server: {
      baseDir: __dirname + '/build'
    },
    middleware: [function (req, res, next) {
      if (
        !_.startsWith(req.url, '/statics') &&
        !_.startsWith(req.url, '/v2')
      ) {
        req.url = '/'
      }

      next();
    }],
  })

  gulp.watch('src/**/*.less', ['less-dev'])
  gulp.watch('index.html', ['html-watch'])
  gulp.watch([
    'src/**/*.js',
    'src/**/*.jsx',
  ], ['webpack:build-dev'])
});

/* Development Mode Settings */
var webpackDevConfig = Object.create(webpackConfig);
webpackDevConfig.debug = true;
var devCompiler = webpack(webpackDevConfig);

gulp.task('webpack:build-dev', function (callback) {
  devCompiler.run(function (err, stats) {
    if (err) {
      throw new gutil.PluginError('webpack:build-dev', err);
    }

    gutil.log('[webpack:build-dev]', stats.toString({
      colors: true,
    }));

    browserSync.reload();
    callback();
  });
});

gulp.task('libjs-dev', function () {
  gulp
    .src([
      'bower_components/tus-js-client/dist/tus.js',
    ])
    .pipe(concat('lib.js'))
    .pipe(gulp.dest('./build/statics/js'));
})

gulp.task('libjs', function () {
  gulp
    .src([
      'bower_components/tus-js-client/dist/tus.js',
    ])
    .pipe(concat('lib.js'))
    .pipe(rename({ suffix: '.min' }))
    .pipe(uglify())
    .pipe(gulp.dest('./build/statics/js'));
})

gulp.task('default', ['build']);
gulp.task('dev', ['build-dev']);
