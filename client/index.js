"use strict"

const ApiUrl = "http://localhost:8080"

const doit = document.querySelector(".do_it")

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
  // const response = await getId()
  // await listen(response.uu_id)
  // console.log("response: ", response);

  // freakingWindow = window.open(response.auth_url, "_blank")

  const popup = window.open(
    'https://t.me/samgasper', // URL для авторизации
    'authPopup',                       // имя окна
    'width=500,height=600,left=100'            // параметры окна
  );
})