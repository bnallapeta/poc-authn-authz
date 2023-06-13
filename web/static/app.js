document.getElementById('getNamespaces').addEventListener('click', function() {
    fetch('/api/getnamespaces')
      .then(response => response.json())
      .then(data => alert(JSON.stringify(data)));
  });
  
  document.getElementById('getPods').addEventListener('click', function() {
    fetch('/api/getpods')
      .then(response => response.json())
      .then(data => alert(JSON.stringify(data)));
  });
  