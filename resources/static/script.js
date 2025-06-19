const resDiv = document.getElementById('results');
const baseUrl = window.location.origin;
const imageUrl = `${baseUrl}/collage`;
const textUrl = `${baseUrl}/altText`;

document.getElementById("collageForm").addEventListener("submit", function(event) {
  event.preventDefault();
  resDiv.innerHTML = '';

  const user = document.getElementById('user').value;
  const size = document.getElementById('size').value;
  const period = document.getElementById('period').value;
  const captions = document.getElementById('captions').checked;
  const plays = document.getElementById('plays').checked;

  const queryParams = new URLSearchParams({
    user,
    size,
    period,
  });

  if (captions) queryParams.append("captions", "on");
  if (plays)    queryParams.append("plays", "on");

  const imageQ = `${imageUrl}?${queryParams.toString()}`;
  const img = document.createElement('img');
  img.src = imageQ;
  img.loading = 'lazy';
  img.id = 'collage';
  resDiv.appendChild(img);

  const textQ  = `${textUrl}?${queryParams.toString()}`;

  fetch(textQ)
    .then(response => {
      if (!response.ok) throw new Error("alt text error");
      return response.text();
    })
    .then(data => {
      document.getElementById('altText')?.remove();
      const altDiv = document.createElement('div');
      altDiv.id = 'altText';

      const altPre = document.createElement('pre');
      altPre.textContent = data;
      altPre.id = 'prealtText';

      const copyBtn = document.createElement('button');
      copyBtn.textContent = 'Copy Alt Text';
      copyBtn.id = 'copyBtn';
      copyBtn.onclick = () => {
        navigator.clipboard.writeText(data)
          .catch(err => console.error("Clipboard copy failed", err));
      };

      altDiv.appendChild(copyBtn);
      altDiv.appendChild(altPre);
      resDiv.appendChild(altDiv);
    })
    .catch(err => console.error(err));
});
