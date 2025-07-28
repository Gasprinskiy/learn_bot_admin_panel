"use strict"

const ApiUrl = "http://localhost:8080"

const doit = document.querySelector(".do_it")
const dot = document.querySelector('.dot')

let freakingWindow = null

async function getId() {
  try {
    const response = await fetch(`${ApiUrl}/auth/temp_data`, {
      method: "GET"
    })

    return response.json()
  } catch (e) {
    console.error(e)
  }
}

async function listen(authID) {
  const url = `${ApiUrl}/auth/listen?auth_id=${encodeURIComponent(authID)}`;
  const source = new EventSource(url);

  source.addEventListener("done", (event) => {
    console.log(`✅ Событие завершено: `, JSON.parse(event.data));
    freakingWindow.close()
    source.close();
  });

  source.addEventListener("error", (event) => {
    console.log(`❌ Событие завершено: `, event);
    freakingWindow.close()
    source.close();
  });
}



doit.addEventListener("click", async () => {
  // const centerX = (window.screen.width - 500) / 2;
  // const centerY = (window.screen.height - 600) / 2;

  const screenWidth = window.screen.availWidth;
  const screenHeight = window.screen.availHeight;
  const windowWidth = 500;
  const windowHeight = 600;

  const centerX = window.screen.availLeft + (screenWidth - windowWidth) / 2;
  const centerY = window.screen.availTop + (screenHeight - windowHeight) / 2;

  // const response = await getId()
  // await listen(response.uu_id)

  freakingWindow = window.open(
    'https://t.me',
    'authPopup',
    `top=${centerY},left=${centerX},width=500,height=600,toolbar=no,menubar=no,resizable=yes`
  );

  // freakingWindow.addEventListener('onbeforeunload', () => {
  //   console.log("close");
  // })

  const interval = setInterval(() => {
    console.log("freakingWindow.closed: ", freakingWindow.closed);
    if (freakingWindow.closed) {
      clearInterval(interval)
    }
  }, 500)
})