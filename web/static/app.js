document.getElementById('getNamespaces').addEventListener('click', function() {
  makeRequest('/api/getnamespaces');
});

document.getElementById('getPods').addEventListener('click', function() {
  makeRequest('/api/getpods');
});

function makeRequest(apiPath) {
  fetch(apiPath)
    .then(response => {
      if (response.status === 401) {
        // Unauthorized, start a new OAuth2 flow
        initiateOAuth2Flow(apiPath);
      } else {
        // Handle the response
        console.log(response);
      }
    })
    .catch(error => console.error('Error:', error));
}

function initiateOAuth2Flow(apiPath) {
  // Redirect the user to the start of a new OAuth2 flow
  const redirectUri = `${window.location.origin}/auth/start?next=${encodeURIComponent(apiPath)}`;
  window.location.href = redirectUri;
}
