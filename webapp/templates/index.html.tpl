<!DOCTYPE html>
<html lang="en">
  <head>
    <link href="/assets/css/bootstrap.css" rel="stylesheet">
    <style type="text/css">
      body {
        padding-top: 60px;
        padding-bottom: 40px;
      }
      .sidebar-nav {
        padding: 9px 0;
      }
    </style>
  </head>
  <body>
    <div class="navbar navbar-fixed-top">
      <div class="navbar-inner">
        <div class="container-fluid">
          <a class="btn btn-navbar" data-toggle="collapse" data-target=".nav-collapse">
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
          </a>
          <a class="brand" href="/">Language Detector</a>
            <div class="nav-collapse">
              <ul class="nav">
                <li class="active"><a href="#">Home</a></li>
                <li><a href="#about">About</a></li>
              </ul>
          </div><!--/.nav-collapse -->
        </div>
      </div>
    </div>
    <div class='container-fluid'>
      <div class="span9"> 
        <p>URL, text... anything</p>
        <div class="row-fluid">
          <form id="form" class="well form-vertical">
            <textarea id="text" name="content" class="field span12"></textarea>
            <button class="btn btn-primary btn-large" type="submit">Detect Language</button>
          </form>
            <div class='log'></div>
        </div>
      </div>
      <div class='raw-fluid pull-right'>
        <div class='span3'>
          <div class="well sidebar-nav">
            <ul class="nav nav-list">
              <li class="nav-header">Supported languages</li>
              {{ range $shortName, $fullName := .SupportedLanguages }}
              <li id=langmarker_{{ $shortName}}><a>{{ $fullName }}</a></li>
              {{ end }}
               <li><a href='/train'>Train</a></li>
             </ul>
          </div>
        </div>
      </div>
    </div>
    <hr>

    <footer>
      <p>&copy; MRO 2012</p>
    </footer>

  <script src="https://code.jquery.com/jquery-4.0.0.min.js"></script>
  <script src="/assets/js/bootstrap.min.js"></script>
  <script language="javascript" type="text/javascript" src="/assets/js/app.js"></script>
  </body>
<html>